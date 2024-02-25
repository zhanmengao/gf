package wssrv

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"github.com/zhanmengao/gf/util"
	fclient "github.com/zhanmengao/gf/wscli"
	"github.com/zhanmengao/gf/wscli/conn"
	"github.com/zhanmengao/gf/wscli/conn/conntyp"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type WsServer struct {
	httpServer     *http.Server
	listenCfg      *cfgtp.ListenConfig
	gatewayURI     string
	gatewayHandler wscli.IHandler
	lock           sync.RWMutex
	handler        map[string]TUrlHandler
	readTimeout    time.Duration
	wsUpgrade      *websocket.Upgrader
	opts           []*fclient.ClientOptions
}

func NewWsServer(gSrv *grpcsrv.Server, lis *cfgtp.ListenConfig) *WsServer {
	srv := &WsServer{
		httpServer: &http.Server{
			WriteTimeout:   0,
			MaxHeaderBytes: 1 << 20,
		},
		listenCfg:   lis,
		handler:     make(map[string]TUrlHandler),
		readTimeout: time.Duration(30) * time.Second,
		wsUpgrade: &websocket.Upgrader{
			HandshakeTimeout: 10 * time.Second,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	srv.httpServer.Handler = srv
	//注册接口
	gensrv.RegisterGatesrvBaseServer(gSrv, srv)
	return srv
}

func (s *WsServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	slog.DebugContext(ctx, "ServeHTTP %s ", request.RequestURI)
	switch request.URL.Path {
	case "/":
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	case s.gatewayURI:
		cli, err := s.HandleClient(ctx, writer, request)
		if err != nil {
			slog.WarnContext(ctx, "WsConn %s HandleClient err = %s", util.GetAddrFromRequest(request), err)
			return
		}
		if err = cli.Run(ctx); err != nil {
			slog.WarnContext(ctx, "fClient %s Run err = %s", cli.GetUID(), err)
			return
		}
	default:
		h, ok := s.getHandler(request.URL.Path)
		if ok {
			resPon, err := h.HandleInit(ctx, request, writer)
			if err != nil {
				slog.WarnContext(ctx, "WsConn %s HandleInit err = %s", util.GetAddrFromRequest(request), err)
				return
			}
			conn, err := s.HandleWebSocket(ctx, writer, request, resPon)
			if err != nil {
				slog.WarnContext(ctx, "WsConn %s HandleWebSocket err = %s", conn.GetRealIP(), err)
				return
			}
			if err = conn.Run(ctx, THandlerPacket(h.HandlePacket)); err != nil {
				slog.WarnContext(ctx, "WsConn %s run err = %s", conn.GetRealIP(), err)
				return
			}
		} else {
			writer.WriteHeader(http.StatusNotFound)
			_, _ = writer.Write([]byte("404 Not Found"))
		}
	}
	slog.InfoContext(ctx, "ServeHTTP %s closed", request.RequestURI)
}

func (s *WsServer) getHandler(urlPath string) (h TUrlHandler, exist bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	h, exist = s.handler[urlPath]
	return
}

func (s *WsServer) HandleClient(ctx context.Context, writer http.ResponseWriter, request *http.Request) (cli *fclient.Client, err error) {
	//建立连接
	wsConn, err := s.wsUpgrade.Upgrade(writer, request, nil)
	if err != nil {
		slog.ErrorContext(ctx, "[WsConnection] Upgrade", err)
		return
	}
	cli = fclient.NewWsClient(ctx, request, wsConn, s.readTimeout, s.gatewayHandler, s.opts...)
	return
}

func (s *WsServer) HandleWebSocket(ctx context.Context, writer http.ResponseWriter, request *http.Request, responseHeader http.Header) (c conntyp.IConnect, err error) {
	//建立连接
	wsConn, err := s.wsUpgrade.Upgrade(writer, request, responseHeader)
	if err != nil {
		slog.ErrorContext(ctx, "[WsConnection] Upgrade", err)
		return
	}
	c = conn.NewWebSocketConn(ctx, request, wsConn, s.readTimeout)
	return
}

func (s *WsServer) processWebSocket(ctx context.Context, conn conntyp.IConnect, h TUrlHandler) {
	var err error
	defer func() {
		h.HandleClose(ctx, conn, err)
	}()
	if err = conn.Run(ctx, THandlerPacket(h.HandlePacket)); err != nil {
		slog.WarnContext(ctx, "WsConn %s run err = %s", conn.GetRealIP(), err)
		return
	}
}
