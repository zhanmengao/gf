package wssrv

import (
	"context"
	"github.com/zhanmengao/gf/wscli/conn/conntyp"
	"net/http"
)

type TUrlHandler interface {
	HandleInit(ctx context.Context, request *http.Request, writer http.ResponseWriter) (responseHeader http.Header, err error)
	HandlePacket(ctx context.Context, conn conntyp.IConnect, rcv []byte)
	HandleClose(ctx context.Context, conn conntyp.IConnect, err error)
}

type THandlerPacket func(ctx context.Context, conn conntyp.IConnect, rcv []byte)

func (f THandlerPacket) HandPacket(ctx context.Context, conn conntyp.IConnect, body []byte) {
	f(ctx, conn, body)
}

type OnlineCallBack func(ctx context.Context, count int)
type ConfigMD5ChangedCallBack func(ctx context.Context, configMD5 string)
