package fclient

import (
	"context"
	"github.com/zhanmengao/gf/proto/go/pb"
)

type BroadcastType int32

const (
	BroadcastUserList = iota
	BroadcastAll
)

type ClientInitType int32

const (
	ClientInitFirstPacket = iota //首个Body发Session
	ClientInitHead               //头部发Session
)

type ClientType int32

const (
	ClientTypeWsV1 = iota
	ClientTypeWsV2
)

const (
	SessionKeyDefault = "token"
)

type THandleHook func(ctx context.Context, cli *Client, pkt *pb.Packet) bool

type ClientOptions struct {
	Type           ClientType
	InitType       ClientInitType
	SessionKeyName string
	HandleHook     THandleHook
}

func NewOption() *ClientOptions {
	return &ClientOptions{
		Type:           ClientTypeWsV2,
		InitType:       ClientInitHead,
		SessionKeyName: SessionKeyDefault,
		HandleHook:     nil,
	}
}

type IHandler interface {

	//HandleInit
	//  @Description: 传递sessionID，业务调用FLS接口
	//  @param cli 客户端类封装
	//  @param sessionID SessionID
	//  @return uid 返回UID
	//  @return err 返回错误，出错时连接被关闭
	HandleInit(ctx context.Context, cli *Client, sessionKey string) (uid, deviceToken string, err error)

	//
	// HandleEncode
	//  @Description:加密器
	//  @param ctx
	//  @param cli
	//  @param pkt 原始报文
	//  @return snd 返回编码后的bytes
	//  @return err 出错会断开连接
	HandleEncode(ctx context.Context, cli *Client, pkt *pb.Packet) (snd []byte, err error)

	//
	// HandleDecode
	//  @Description: 解码器
	//  @param ctx
	//  @param rcv 收到原文
	//  @return pkt 返回解出来的packet
	//  @return err 出错会断开连接
	//
	HandleDecode(ctx context.Context, cli *Client, rcv []byte) (pkt *pb.Packet, isSendToSrv bool, err error)

	//HandleClose
	//  @Description: 关闭前的回调，可以发送遗言
	//  @param ctx：上下文
	//  @param cli：客户端类的封装
	//
	HandleClose(ctx context.Context, cli *Client, err error)
}

type IV11Handler interface {
	//HandleInit
	//  @Description: 传递sessionID，业务调用FLS接口
	//  @param cli 客户端类封装
	//  @param sessionID SessionID
	//  @return uid 返回UID
	//  @return err 返回错误，出错时连接被关闭
	HandleInit(ctx context.Context, cli *Client, sessionKey string, rcv []byte) (uid, aesKey string, err error)

	//HandlePacket
	//  @Description: pkt过滤器
	//  @param cli 客户端类封装
	//  @param pkt 消息包
	//  @return ok 是否继续转发
	//
	HandlePacket(ctx context.Context, cli *Client, pkt *pb.Packet) (isSendToSrv bool)

	//HandleClose
	//  @Description: 关闭前的回调，可以发送遗言
	//  @param ctx：上下文
	//  @param cli：客户端类的封装
	//
	HandleClose(ctx context.Context, cli *Client, err error)
}
