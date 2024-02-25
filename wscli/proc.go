package fclient

import (
	"container/list"
	"context"
	ws "github.com/gorilla/websocket"
	framework "github.com/zhanmengao/gf/context"
	"github.com/zhanmengao/gf/proto/go/gerror"
	"github.com/zhanmengao/gf/proto/go/pb"
	"github.com/zhanmengao/gf/wscli/conn"
	"github.com/zhanmengao/gf/wscli/conn/conntyp"
	"log/slog"
	"net/http"
	"time"
)

func NewWsClient(ctx context.Context, request *http.Request, c *ws.Conn, readTimeout time.Duration, handler IHandler, opts ...*ClientOptions) *Client {
	cli := &Client{
		handler:        handler,
		rspPacketCache: list.New(),
		reqPacketCache: make(map[int64]struct{}),
	}
	if len(opts) == 0 {
		opts = append(opts, NewOption())
	}
	for _, opt := range opts {
		cli.option = opt
		switch opt.Type {
		case ClientTypeWsV1:
			cli.conn = conn.NewWebSocketConn(ctx, request, c, readTimeout)
		case ClientTypeWsV2:
			cli.conn = conn.NewWebSocketConnV2(ctx, request, c, readTimeout)
		}
		cli.handleHook = append(cli.handleHook, opt.HandleHook)
	}
	return cli
}

func NewWsClientByConn(ctx context.Context, conn conntyp.IConnect, handler IHandler, opts ...*ClientOptions) *Client {
	cli := &Client{
		conn:           conn,
		handler:        handler,
		rspPacketCache: list.New(),
		reqPacketCache: make(map[int64]struct{}),
	}
	if len(opts) == 0 {
		opts = append(opts, NewOption())
	}
	for _, opt := range opts {
		cli.option = opt
	}
	return cli
}

func GetClient(uid string) (cli *Client) {
	sess, ok := framework.SessPool.GetUserSession(uid)
	if !ok {
		return
	}
	cli = GetClientFromSess(sess)
	return
}

func GetClientFromSess(sess *framework.Session) *Client {
	iClient := sess.GetSessValue(framework.SessionKeyClient)
	if iClient == nil {
		return nil
	}
	cli, ok := iClient.(*Client)
	if !ok {
		return nil
	}
	return cli
}

// BroadcastMsgToUser 网关调用，发消息到uid列表
func BroadcastMsgToUser(ctx context.Context, head *pb.PacketHead, msg framework.IProto, broadcastType BroadcastType, uidList ...string) (resUIDList []string, err error) {
	data, err := msg.Marshal()
	if err != nil {
		err = gerror.ErrServerEncode().SetBasicErr(err)
		return
	}
	pkt := &pb.Packet{
		Head: head,
		Body: data,
	}
	resUIDList = BroadcastPktToUser(ctx, pkt, broadcastType, uidList...)
	return
}

// BroadcastPktToUser 网关调用，发消息到uid列表
func BroadcastPktToUser(ctx context.Context, pkt *pb.Packet, broadcastType BroadcastType, uidList ...string) (resUIDList []string) {
	if broadcastType == BroadcastAll {
		framework.SessPool.ForeachUID(func(uid string, sess *framework.Session) bool {
			if cli := GetClientFromSess(sess); cli != nil {
				if err := cli.SendPacket(ctx, pkt); err != nil {
					slog.WarnContext(ctx, "send [%d] to client [%s]. error = %s ", pkt.Head.Cmd, uid, err.Error())
				} else {
					resUIDList = append(resUIDList, uid)
				}
			}
			return true
		})
	} else if broadcastType == BroadcastUserList {
		for _, uid := range uidList {
			if cli := GetClient(uid); cli != nil {
				if err := cli.SendPacket(ctx, pkt); err != nil {
					slog.WarnContext(ctx, "send [%d] to client [%s]. error = %s ", pkt.Head.Cmd, uid, err.Error())
				} else {
					resUIDList = append(resUIDList, uid)
				}
			}
		}
	}
	return
}

type RangUserCallBack func(ctx context.Context, cli *Client)

func RangeUser(ctx context.Context, cb RangUserCallBack) {
	framework.SessPool.ForeachUID(func(s string, session *framework.Session) bool {
		if cli := GetClientFromSess(session); cli != nil {
			cb(ctx, cli)
		}
		return true
	})
}
