package framework

import (
	"context"
	"github.com/zhanmengao/gf/basedef"
	"github.com/zhanmengao/gf/proto/go/pb"
	"go.opentelemetry.io/otel/propagation"
	"sync"
	"time"
)

const (
	scopeUser           = "user"
	KeyNameUID          = basedef.KeyNameUID
	KeyNameRoute        = basedef.KeyNameRoute
	KeyNameError        = basedef.KeyNameError
	KeyNameRequest      = basedef.KeyNameRequest
	KeyNameResponse     = basedef.KeyNameResponse
	GinKeyFormat        = basedef.GinKeyFormat
	GinKeyBody          = basedef.GinKeyBody
	GinKeyWriteResponse = basedef.GinKeyWriteResponse
	KeyNamePlatform     = basedef.KeyNamePlatform
	ContextBeginKey     = basedef.ContextBeginKey
	KeyNameContext      = basedef.KeyNameContext
	GinKeyContext       = basedef.GinKeyContext
	KeyNameBaseSpanCtx  = basedef.KeyNameBaseSpanCtx
	GinKeyCurrentCtx    = basedef.GinKeyCurrentCtx
	KeyNameFullMethod   = basedef.KeyNameFullMethod
	KeyNameUseTrace     = basedef.KeyNameUseTrace
	KeyNameReqID        = basedef.KeyNameReqID
	KeyNameRemoteIP     = basedef.KeyNameRemoteIP
	KeyNameRemoteSrv    = basedef.KeyNameRemoteSrv
	KeyNameService      = basedef.KeyNameService
	CacheKeyPrefix      = basedef.CacheKeyPrefix
)

type TValWrapper struct {
	Key string
	Val interface{}
}

type Context struct {
	*pb.PacketHead
	sync.RWMutex
	fromAddr    string
	fromSrvName string
	sess        *Session
	data        []TValWrapper
	rspSent     bool //是否已经给前端回过包
	ctx         context.Context
}

func NewContext(c context.Context) *Context {
	ctx := &Context{
		PacketHead: &pb.PacketHead{},
		ctx:        c,
	}
	ctx.SetContextValue(KeyNameContext, ctx)
	return ctx
}

func (c *Context) SetPacketHead(head *pb.PacketHead) {
	c.PacketHead = head
	c.SetContextValue(KeyNameUID, head.UID)
	c.SetContextValue(KeyNameRoute, head.RKey)

}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if c.ctx != nil {
		return c.ctx.Deadline()
	}
	return
}

func (c *Context) Done() <-chan struct{} {
	if c.ctx != nil {
		return c.ctx.Done()
	}
	return nil
}

func (c *Context) Err() error {
	if c.ctx != nil {
		return c.ctx.Err()
	}
	return nil
}

func (c *Context) Value(iKey interface{}) interface{} {
	key, ok := iKey.(string)
	if ok {
		if reply := c.GetContextValue(key); reply != nil {
			return reply
		}
	}
	if c.ctx != nil {
		return c.ctx.Value(iKey)
	}
	return nil
}

func (c *Context) GetContentValueList() []TValWrapper {
	return c.data
}

// GetMetaDataValue 获取Meta Value
func (c *Context) GetMetaDataValue(key string) (value string, exist bool) {
	c.RLock()
	defer c.RUnlock()
	value, exist = c.MetaData[key]
	return
}

// SetMetaDataValue 设置Meta Value
func (c *Context) SetMetaDataValue(key, value string) {
	c.Lock()
	defer c.Unlock()
	if c.MetaData == nil {
		c.MetaData = map[string]string{}
	}
	//先找有没有
	c.MetaData[key] = value
	c.setContextValue(key, value)
}

func (c *Context) GetMapCarrierCopy() propagation.MapCarrier {
	c.RLock()
	defer c.RUnlock()
	md := c.GetMetaData()
	if md == nil {
		return nil
	}
	ret := propagation.MapCarrier(make(map[string]string, len(md)))
	for k, v := range md {
		ret[k] = v
	}
	return ret
}

func (c *Context) GetContextValue(key string) (value interface{}) {
	c.RLock()
	defer c.RUnlock()
	return c.getContextValue(key)
}

func (c *Context) getContextValue(key string) (value interface{}) {
	for _, v := range c.data {
		if v.Key == key {
			return v.Val
		}
	}
	return nil
}

func (c *Context) SetContextValue(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.setContextValue(key, value)
}

func (c *Context) setContextValue(key string, value interface{}) {
	i := 0
	for ; i < len(c.data); i++ {
		if c.data[i].Key == key {
			c.data[i].Val = value
			return
		}
	}
	c.data = append(c.data, TValWrapper{
		key, value,
	})
}

func (c *Context) GetAddr() string {
	return c.fromAddr
}

func (c *Context) GetSrvName() string {
	return c.fromSrvName
}

func (c *Context) GetSession() *Session {
	return c.sess
}

func (c *Context) SetSession(s *Session) {
	c.sess = s
}

func (c *Context) SetRspSent(rspSent bool) {
	c.rspSent = rspSent
}

func (c *Context) GetRspSent() bool {
	return c.rspSent
}

func (c *Context) SetFromSrvName(srvName string) {
	c.fromSrvName = srvName
}

func (c *Context) SetFromSrvAddr(addr string) {
	c.fromAddr = addr
}

func (c *Context) GetFromSrvName() (srvName string) {
	srvName = c.fromSrvName
	return
}

func (c *Context) GetFromSrvAddr() (addr string) {
	addr = c.fromAddr
	return
}
