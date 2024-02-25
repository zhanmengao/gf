package fclient

import (
	"container/list"
	"context"
	uuid "github.com/satori/go.uuid"
	framework "github.com/zhanmengao/gf/context"
	"github.com/zhanmengao/gf/proto/go/gerror"
	"github.com/zhanmengao/gf/proto/go/pb"
	"github.com/zhanmengao/gf/trace"
	"github.com/zhanmengao/gf/trace/tracetyp"
	"github.com/zhanmengao/gf/wscli/conn/conntyp"
	"github.com/zhanmengao/gf/xtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	conn       conntyp.IConnect
	Sess       *framework.Session
	handler    IHandler
	isInit     bool
	isClosed   int32
	lastSeq    int32
	option     *ClientOptions
	sessionKey string
	//包缓存
	lock              sync.RWMutex
	dbSession         *pb.DBSessionKey
	pktOpts           int64 //最后一个收到的Opts
	lastSaveDBSession time.Time
	handleHook        []THandleHook
	rspPacketCache    *list.List
	reqPacketCache    map[int64]struct{}
}

func (c *Client) Run(ctx context.Context) (err error) {
	defer func() {
		c.CloseNotWait(ctx, err)
	}()
	//未初始化，先读sessionID
	if c.option.InitType == ClientInitHead {
		c.sessionKey = getFromHeader(c.conn.GetRequest(), c.option.SessionKeyName)
		if err = c.initClient(ctx); err != nil {
			return
		}
	}
	if err = c.conn.Run(ctx, c); err != nil {
		return
	}
	return
}

// 从协议头里获得参数
func getFromHeader(req *http.Request, key string) string {
	param := req.Header.Get(key)
	if param == "" {
		param = req.URL.Query().Get(key)
	}
	return param
}

func (c *Client) initClient(ctx context.Context) (err error) {
	//先看SessionKey是否存在，存在的话直接算登录成功
	var ok bool
	c.dbSession, ok, err = ConnectDB.GetDBSessionKey(ctx, c.sessionKey)
	if err != nil {
		return
	}
	var uid string
	var deviceToken = c.dbSession.DeviceToken
	//读到了，说明是重连
	if ok {
		uid = c.dbSession.UID
	} else {
		if uid, deviceToken, err = c.handler.HandleInit(ctx, c, c.sessionKey); err != nil {
			return
		}
		c.dbSession.UID = uid
		//缓存6小时
		if err = c.SaveDBSession(ctx); err != nil {
			return
		}
	}
	for _, pkt := range c.dbSession.PacketList {
		c.rspPacketCache.PushBack(pkt)
	}
	if uid == "" {
		slog.WarnContext(ctx, "get a nil uid .from %s .create at %s .ConnId = %d ", c.conn.GetRealIP(),
			time.Unix(c.conn.GetCreateTime(), 0).String(), c.conn.GetConnID())
		err = gerror.ErrServerInternalUnknown().Format("uid is empty")
		return
	}
	c.isInit = true
	//读DB，看是否有老连接
	var dbConn *pb.DBConnect
	if dbConn, ok, err = ConnectDB.GetDBConnect(ctx, uid); err != nil {
		return
	}
	if ok {
		//调用其他服务踢人
		_, err = GateCli.KickUserToAddr(ctx, &pb.GateKickUserReq{
			UID:         uid,
			DeviceToken: deviceToken,
		}, uid, dbConn.GatewayAddr)
		if err != nil {
			slog.ErrorContext(ctx, "KickUserToAddr uid[%s] device[%s] error[%s]", uid, deviceToken, err)
		} else {
			slog.InfoContext(ctx, "KickUserToAddr uid[%s] device[%s] success", uid, deviceToken)
		}
	}

	if ctx, err = c.SetUID(ctx, uid); err != nil {
		return
	}
	//写DB
	dbConn.GatewayAddr = fgrpc.GetLocalAddr()
	if err = ConnectDB.SetDBConnect(ctx, dbConn); err != nil {
		return
	}
	return
}

func (c *Client) HandPacket(ctx context.Context, conn conntyp.IConnect, rcv []byte) {
	//出现解包等错误，直接关连接
	var err error
	var cmd int32
	defer func() {
		if err != nil {
			c.CloseNotWait(ctx, err)
		}
	}()
	if len(rcv) <= 0 {
		return
	}

	pktContext := framework.NewContext(ctx)
	ctx = pktContext

	if !c.isInit {
		c.sessionKey = string(c.conn.GetFirstBody())
		if err = c.initClient(ctx); err != nil {
			return
		}
	}
	if c.Sess == nil {
		c.Sess, ctx, err = framework.SessPool.NewUserSession(ctx, c.GetUID())
	}
	pktContext.SetSession(c.Sess)

	//回调
	pkt, sendToSrv, err := c.handler.HandleDecode(ctx, c, rcv)
	if err != nil {
		return
	}

	if pkt.Head != nil {
		c.pktOpts = pkt.Head.Opts
	}
	c.Sess.UpdateActiveTime()
	//触发basecmd检测
	isBaseCMD := c.handleBaseCmd(ctx, pkt)

	//收到包了，检测AckPeerID，看是否有对端没收到的包
	//处理一下包缓存
	needHandle := c.checkPeerRcv(ctx, pkt)
	//回调如果说这个包不发，则直接返回
	if !sendToSrv || !needHandle || isBaseCMD {
		return
	}
	c.lastSeq = pkt.Head.Seq
	//每次收包会恢复可能丢失的session
	if ctx, err = c.SetUID(ctx, c.GetUID()); err != nil {
		return
	}
	//为空则生成一个ReqId
	if pkt.Head.ReqId == "" {
		pkt.Head.ReqId = uuid.NewV5(uuid.NewV4(), "req").String()
	}
	ctx = framework.SetMetaValue(ctx, framework.KeyNameReqID, pkt.Head.ReqId)
	pktContext.PacketHead = pkt.Head
	pktContext.SetContextValue(framework.KeyNameRemoteIP, c.GetConnect().GetRealIP())
	ok := true
	for _, cb := range c.handleHook {
		if cb != nil {
			if ok = cb(ctx, c, pkt); !ok {
				break
			}
		}
	}
	//如果hook终止了转发，那么就不再继续了
	if !ok {
		return
	}
	var srv, addr string
	//是否触发trace：网关层不做判断，由srv层做判断。server层能拿到接口名，网关只有CMD
	var span *tracetyp.RpcSpan
	ctx, span = trace.NewRpcSpanByUse(ctx, "frpc", "", tracetyp.RpcOpSend, false)
	defer func() {
		span.SetBodySize(len(pkt.Body))
		span.SetUID(c.GetUID())
		span.SetIsLocal(false)
		span.SetRemoteSrv(srv)
		span.SetRemoteAddr(addr)
		span.End(pkt.Head, nil, err)
	}()
	if srv, addr, err = c.SendToSrv(ctx, pkt); err != nil {
		xlog.Errorf(ctx, "SendToSrv[%s~%d]: %s", c.GetUID(), pkt.Head.Cmd, err.Error())
	}
	return
}

func (c *Client) SendToSrv(ctx context.Context, pkt *pb.Packet) (srv, addr string, err error) {
	var ok bool
	//寻找cmd对srv的映射
	if srv, err = frpc.CMDToSrvName(pkt.Head.Cmd); err != nil {
		if err != nil {
			xlog.Errorf(ctx, "CMDToSrvName[%s~%d]: %s", c.GetUID(), pkt.Head.Cmd, err.Error())
			return
		}
	} else {
		//是否有固定Route
		addr, ok = c.GetRouteAddr(srv)
		if ok && addr != "" {
			//尝试定向。成功了就直接返回
			if err = frpc.SendPktToAddr(ctx, srv, addr, pkt); err == nil {
				return
			}
		}
		//如果节点不在了，就重新Route
		addr, err = frpc.SendPktToSrv(ctx, srv, c.GetUID(), pkt.Head.RKey, pkt.Head.Cmd, pkt)
		if err != nil {
			return
		} else {
			//设置route地址缓存
			c.SetRouteAddr(srv, addr)
		}
	}
	return
}

// SendPacket 发送通信格式的数据
func (c *Client) SendPacket(ctx context.Context, pkt *pb.Packet) (err error) {
	mc := pkt.GetHead().MetaData
	if len(mc) > 0 {
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(pkt.GetHead().MetaData))
	}
	//给SendID、AckPeerID赋值
	c.dbSession.SendRspID++
	pkt.Head.SendID = c.dbSession.SendRspID
	//每次回包，做包缓存
	c.cachePacket(ctx, pkt)

	return c.sendPacket(ctx, pkt)
}

func (c *Client) sendPacket(ctx context.Context, pkt *pb.Packet) (err error) {
	pkt.Head.Time = xtime.Millisecond()
	//返回当前接收到的最大请求ID
	pkt.Head.AckPeerID = c.dbSession.AckPeerID
	//create
	var span *tracetyp.RpcSpan
	ctx, span = ftrace.NewRpcSpanByUse(ctx, "frpc", framework.GetCMDNameMapping(strconv.Itoa(int(pkt.Head.Cmd-1))), tracetyp.RpcOpSend, pkt.Head.UseTrace)
	if span.IsValid() {
		span.SetBodySize(len(pkt.Body))
		span.SetRemoteAddr(c.GetConnect().GetRealIP())
		span.SetUID(c.GetUID())
		span.SetIsLocal(false)
		defer func() {
			span.End(pkt.Head, nil, err)
		}()
	}
	//直接回调
	sendBuf, err := c.handler.HandleEncode(ctx, c, pkt)
	if err != nil {
		return
	}
	err = c.conn.Write(ctx, sendBuf)

	if err != nil {
		return
	}
	return
}

// Close 关闭client并等待关闭事件完成。必须在Handler外调用，在Handler内调用会导致阻塞。
func (c *Client) Close(ctx context.Context, err error) error {
	c.CloseNotWait(ctx, err)
	//这里会阻塞等client关闭完成
	select {
	case <-c.conn.GetClosedChannel():
		return nil
	case <-time.After(time.Second * time.Duration(5)):
		err = xerrors.ErrServerNetwork().Format("closed %s time out", c.GetUID())
		return err
	}
}

// CloseNotWait 关闭cli且不等待结果。
func (c *Client) CloseNotWait(ctx context.Context, err error) {
	if atomic.CompareAndSwapInt32(&c.isClosed, 0, 1) {
		//关闭前，触发遗言
		c.handler.HandleClose(ctx, c, err)
		//实际关闭连接
		c.conn.Close(ctx)
		//删除会话
		framework.SessPool.DeleteUserSession(ctx, c.GetUID())
		pressure.CloseUID(c.GetUID())
		//删除连接
		if _, err = ConnectDB.DeleteDBConnect(ctx, c.GetUID()); err != nil {
			xlog.Errorf(ctx, "DeleteDBConnect [%s] error = %s", c.GetUID(), err)
		}
		//存储Session
		if err = c.SaveDBSession(ctx); err != nil {
			xlog.Errorf(ctx, "SaveDBSession [%s] error = %s", c.GetUID(), err)
		}
	}
}

func (c *Client) SaveDBSession(ctx context.Context) (err error) {
	c.dbSession.PacketList = c.getAllPacket(ctx)
	if err = ConnectDB.SetEXDBSessionKey(ctx, c.dbSession, 6*60*60); err != nil {
		return
	}
	c.lastSaveDBSession = time.Now()
	return
}

func (c *Client) cachePacket(ctx context.Context, packet *pb.Packet) {
	//如果未启用包缓存，直接return
	if !baseutil.GetBit(int32(c.pktOpts), int(pb.PACKET_HEAD_OPT_SAFE_NETWORK_OPEN)) {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.rspPacketCache.PushBack(packet)
	////每间隔进行一次落地DB
	//if time.Since(c.lastSaveDBSession) > time.Duration(10)*time.Second {
	//	if err := c.SaveDBSession(ctx); err != nil {
	//		xlog.Errorf(ctx, "SaveDBSession error [%s]", err)
	//	}
	//	//成不成功都记录时间，避免持续压力DB
	//	c.lastSaveDBSession = time.Now()
	//}
}

func (c *Client) handleBaseCmd(ctx context.Context, pkt *pb.Packet) (isBaseCMD bool) {
	if pkt.Head.Cmd > int32(pb.BASE_CMD_END) {
		return
	}
	isBaseCMD = true
	switch pkt.Head.Cmd {
	case int32(pb.BASE_CMD_RECONNECT_REQ):
		//处理重连
		_, _ = c.getCachePacket(ctx, pkt)
		//截一下，直接把所有包发回去
		rsp := &pb.ReConnRsp{
			Head:    &pb.RspHead{},
			RspList: c.getAllPacket(ctx),
		}
		body, _ := rsp.Marshal()
		rspPkt := &pb.Packet{
			Head: &pb.PacketHead{
				Cmd:            int32(pb.BASE_CMD_RECONNECT_RSP),
				Seq:            0,
				Time:           xtime.Millisecond(),
				Opts:           pkt.Head.Opts,
				RKey:           "",
				ClientVer:      "",
				ClientIP:       "",
				Mod:            0,
				ReqId:          pkt.Head.ReqId,
				Event:          false,
				ClientRevision: "",
				FlowID:         0,
				Reserved:       nil,
				UID:            c.GetUID(),
				TableID:        0,
				MetaData:       pkt.Head.MetaData,
				UseTrace:       pkt.Head.UseTrace,
				RoleID:         pkt.Head.RoleID,
				SendID:         c.dbSession.SendRspID,
				AckPeerID:      c.dbSession.AckPeerID,
			},
			Body: body,
		}
		err := c.sendPacket(ctx, rspPkt)
		if err != nil {
			xlog.Warnf(ctx, "send basecmd[%d] to uid[%s] error[%s]", pkt.Head.Cmd, c.GetUID(), err)
		}
	}
	return
}

func (c *Client) checkPeerRcv(ctx context.Context, pkt *pb.Packet) (needHandle bool) {
	//如果未启用包缓存，直接return
	if !baseutil.GetBit(int32(c.pktOpts), int(pb.PACKET_HEAD_OPT_SAFE_NETWORK_OPEN)) {
		needHandle = true
		return
	}
	pktList, needHandle := c.getCachePacket(ctx, pkt)
	for _, cPkt := range pktList {
		if err := c.sendPacket(ctx, cPkt); err != nil {
			xlog.Errorf(ctx, "sendPacket cmd[%d] uid[%s] error", cPkt.Head.Cmd, c.GetUID())
		}
	}
	return
}

func (c *Client) getCachePacket(ctx context.Context, pkt *pb.Packet) (pktList []*pb.Packet, needHandle bool) {
	func() {
		if c.dbSession.AckPeerID < pkt.Head.SendID {
			c.dbSession.AckPeerID = pkt.Head.SendID
		}
		needHandle = c.SplitCachePacket(pkt)
	}()

	//看哪些包要进行重传：回包超过500ms未ack的
	now := xtime.Millisecond()
	c.lock.RLock()
	defer c.lock.RUnlock()
	for element := c.rspPacketCache.Front(); element != nil; element = element.Next() {
		cPkt := element.Value.(*pb.Packet)
		//检测该包是否已经处理过，处理过了直接回结果
		if pkt.Head.Seq+1 == cPkt.Head.Seq {
			pktList = append(pktList, cPkt)
			needHandle = false
		} else if now-cPkt.Head.Time > 500 {
			//回包超过500ms未ack的
			pktList = append(pktList, cPkt)
		}
	}
	return
}

func (c *Client) getAllPacket(ctx context.Context) (pktList []*pb.Packet) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for element := c.rspPacketCache.Front(); element != nil; element = element.Next() {
		pktList = append(pktList, element.Value.(*pb.Packet))
	}
	return
}

func (c *Client) SplitCachePacket(pkt *pb.Packet) (needHandle bool) {
	needHandle = true
	c.lock.Lock()
	defer c.lock.Unlock()
	//检测当前确认号
	for element := c.rspPacketCache.Front(); element != nil; {
		cPkt := element.Value.(*pb.Packet)
		//把未接收的包remove
		if cPkt.Head.SendID <= pkt.Head.AckPeerID {
			next := element.Next()
			c.rspPacketCache.Remove(element)
			element = next
		} else {
			break
		}
	}
	//检测该包是否已经处理过，处理过了不再处理
	if _, ok := c.reqPacketCache[pkt.Head.SendID]; ok {
		needHandle = false
	}
	c.reqPacketCache[pkt.Head.SendID] = struct{}{}
	return
}
