package wsv2

import (
	"context"
	ws "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/zhanmengao/gf/util"
	"github.com/zhanmengao/gf/xtime"
	"log/slog"
	"net/http"
	"time"
)

func NewWSConn(ctx context.Context, request *http.Request, ws *ws.Conn, readTimeout time.Duration) *Conn {
	conn := &Conn{
		conn:        ws,
		readTimeout: readTimeout,
		closedCh:    make(chan struct{}, 1),
		connID:      uuid.NewV5(uuid.NewV4(), "ws").String(),
		createTime:  xtime.Unix(),
		closeTime:   0,
		request:     request,
	}
	if ws != nil {
		conn.remoteAddr = ws.RemoteAddr()
		conn.localAddr = ws.LocalAddr()
	}
	//读X-Forward-For
	if request != nil {
		conn.forwardIP = util.GetAddrFromRequest(request)
	}
	return conn
}

func (c *Conn) getConn() *ws.Conn {
	return c.conn
}

func (c *Conn) GetIP() string {
	return c.forwardIP
}

func (c *Conn) GetConnID() string {
	return c.connID
}

func (c *Conn) GetCreateTime() int64 {
	return c.createTime
}

func (c *Conn) Write(ctx context.Context, body []byte) (err error) {
	err = sendQueue.PushJob(ctx, c.connID, func(ctx context.Context) {
		err2 := c.write(ctx, body)
		if err2 != nil {
			c.Close(ctx)
			slog.WarnContext(ctx, "%s conn close ", c.GetIP())
		}
	})
	return
}

func (c *Conn) Close(ctx context.Context) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_ = c.conn.Close()
}
