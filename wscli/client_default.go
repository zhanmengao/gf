package fclient

import (
	"context"
	"github.com/zhanmengao/gf/proto/go/pb"
	"log/slog"
)

type DefaultPlainHandle struct {
}

func (*DefaultPlainHandle) HandleEncode(ctx context.Context, cli *Client, pkt *pb.Packet) (snd []byte, err error) {
	snd, err = pkt.Marshal()
	return
}

func (*DefaultPlainHandle) HandleDecode(ctx context.Context, cli *Client, rcv []byte) (pkt *pb.Packet, isSendToSrv bool, err error) {
	pkt = &pb.Packet{}
	isSendToSrv = true
	err = pkt.Unmarshal(rcv)
	return
}

func (*DefaultPlainHandle) HandleClose(ctx context.Context, cli *Client, err error) {
	slog.InfoContext(ctx, "cli[%s] was closed.error[%s]", cli.GetUID(), err)
}
