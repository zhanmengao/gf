package framework

import "github.com/zhanmengao/gf/proto/go/pb"

"

type IProto interface {
	// Message
	Reset()
	String() string
	ProtoMessage()
	// Marshaler
	Marshal() ([]byte, error)
	// Unmarshaler
	Unmarshal([]byte) error
	Size() int
}

type IReply interface {
	IProto
	GetHead() *pb.RspHead
}
