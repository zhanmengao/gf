package wssrv

import (
	"context"
	"github.com/zhanmengao/gf/proto/go/gerror"
	"github.com/zhanmengao/gf/proto/go/pb"
	fclient "github.com/zhanmengao/gf/wscli"
	"github.com/zhanmengao/gf/xtime"
	"log/slog"
)

func (s *WsServer) KickUser(ctx context.Context, req *pb.GateKickUserReq) (rsp *pb.CustomErrorRsp, err error) {
	rsp = &pb.CustomErrorRsp{}
	//关闭老连接
	if old := fclient.GetClient(req.GetUID()); old != nil {
		if err = old.Close(ctx, gerror.ErrServerNetwork().Format("del old client")); err != nil {
			return
		}
	}
	return
}

func (s *WsServer) NotifyToUser(ctx context.Context, req *pb.NotifyToUserReq) (rsp *pb.NotifyToUserRsp, err error) {
	rsp = &pb.NotifyToUserRsp{}
	if req.Broadcast {
		//全员广播,找到所有玩家
		fclient.RangeUser(ctx, func(ctx context.Context, cli *fclient.Client) {
			if err2 := s.sendPacketToUser(ctx, cli, req.CMD, req.Data); err2 != nil {
				slog.ErrorContext(ctx, "sendPacketToUser uid[%s].cmd[%d],err[%s]", cli.GetUID(), req.CMD, err2)
			}
		})
	} else {
		for _, uid := range req.UIDList {
			cli := fclient.GetClient(uid)
			if cli != nil {
				if err2 := s.sendPacketToUser(ctx, cli, req.CMD, req.Data); err2 != nil {
					slog.ErrorContext(ctx, "sendPacketToUser uid[%s].cmd[%d],err[%s]", cli.GetUID(), req.CMD, err2)
				}
			}
		}
	}
	return
}

func (s *WsServer) sendPacketToUser(ctx context.Context, cli *fclient.Client, cmd int32, data []byte) (err error) {
	pkt := &pb.Packet{
		Head: &pb.PacketHead{
			Cmd:            cmd,
			Seq:            0,
			Time:           xtime.Millisecond(),
			Opts:           0,
			RKey:           "",
			ClientVer:      "",
			ClientIP:       cli.GetConnect().GetRealIP(),
			Mod:            0,
			ReqId:          "",
			Event:          false,
			ClientRevision: "",
			FlowID:         0,
			Reserved:       nil,
			UID:            cli.GetUID(),
			TableID:        0,
			MetaData:       nil,
			UseTrace:       false,
			RoleID:         "",
			SendID:         0,
			AckPeerID:      0,
		},
		Body: data,
	}
	if err = cli.SendPacket(ctx, pkt); err != nil {
		return
	}
	return
}

func (s *WsServer) handleSrvPacket(ctx context.Context, srvName, addr string, data []byte) (err error) {
	//将内网RPC包，转为外网包返回
	if len(data) <= 0 {
		err = gerror.ErrServerBadParam().Format("data is nil")
		return
	}
	pkt := &pb.Packet{}
	if err = pkt.Unmarshal(data); err != nil {
		err = gerror.ErrServerDecode().SetBasicErr(err)
		return
	}
	cli := fclient.GetClient(pkt.Head.UID)
	if cli == nil {
		return
	}
	if err = cli.SendPacket(ctx, pkt); err != nil {
		slog.ErrorContext(ctx, "send [%d] to [%s] .err = %s ", pkt.Head.Cmd, cli.GetUID(), err.Error())
	} else {
		slog.DebugContext(ctx, "send [%d] to [%s] success ", pkt.Head.Cmd, cli.GetUID())
	}
	return
}
