package conn

import (
	"context"
	ws "github.com/gorilla/websocket"
	"github.com/zhanmengao/gf/wscli/conn/conntyp"
	"github.com/zhanmengao/gf/wscli/conn/internal/wsv2"
	"net/http"
	"time"
)

func NewWebSocketConnV2(ctx context.Context, request *http.Request, conn *ws.Conn, readTimeout time.Duration) conntyp.IConnect {
	return wsv2.NewWSConn(ctx, request, conn, readTimeout)
}
