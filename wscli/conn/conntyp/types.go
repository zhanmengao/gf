package conntyp

import (
	"context"
	"net"
	"net/http"
)

type IConnect interface {
	//Run
	//  @Description: 连接开始接收事件，只可调用一次
	//  @param ctx 上下文信息
	//  @param handler 注入的处理函数
	//  @return err 连接终止时返回error
	Run(ctx context.Context, handler IHandler) (err error)
	//Write
	//  @Description: 向对端连接回包
	//  @param ctx：上下文
	//  @param body：发送的消息
	//  @return err： 错误
	Write(ctx context.Context, body []byte) (err error)

	GetRealIP() string
	GetRemoteAddr() net.Addr
	GetLocalAddr() net.Addr
	GetCreateTime() int64
	GetConnID() string
	GetRequest() *http.Request
	GetFirstBody() []byte
	//Close
	//  @Description: 立刻关闭连接，可以在Run出现错误后调用
	Close(ctx context.Context)

	GetClosedChannel() chan struct{}
}

type IHandler interface {
	HandPacket(ctx context.Context, conn IConnect, body []byte)
}
