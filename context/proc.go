package framework

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zhanmengao/gf/basedef"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"github.com/zhanmengao/gf/trace/tracetyp"
	"github.com/zhanmengao/gf/util"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var SessPool = NewSessionPool(context.Background())

var BasicConfig cfgtp.IBasicConfig

func GetContext(c context.Context) (ctx *Context, ok bool) {
	if c == nil {
		return
	}
	iContext := c.Value(KeyNameContext)
	if ctx, ok = iContext.(*Context); ok {
		return
	}
	return
}

func GetGinContext(c context.Context) (ctx *gin.Context, ok bool) {
	if c == nil {
		return
	}
	iContext := c.Value(GinKeyContext)
	if ctx, ok = iContext.(*gin.Context); ok {
		return
	}
	return
}

func SetContextValue(c context.Context, key string, value interface{}) bool {
	if ctx, ok := GetContext(c); ok {
		ctx.SetContextValue(key, value)
		return true
	}
	if ctx, ok := GetGinContext(c); ok {
		ctx.Set(key, value)
		return true
	}
	return false
}

func GetMetaValue(c context.Context, key string) (value string, exist bool) {
	//frpc
	if ctx, ok := GetContext(c); ok {
		if value, exist = ctx.GetMetaDataValue(key); exist {
			return
		}
	}
	//grpc
	md, ok := metadata.FromIncomingContext(c)
	if ok {
		sl := md.Get(key)
		if len(sl) > 0 {
			value = sl[0]
			exist = true
		}
	}
	return
}

func SetMetaValue(c context.Context, key, value string) context.Context {
	//frpc
	if ctx, ok := GetContext(c); ok {
		ctx.SetMetaDataValue(key, value)
	}
	//grpc
	if requestMetadata, ok := metadata.FromOutgoingContext(c); ok {
		requestMetadata.Set(key, value)
		c = metadata.NewOutgoingContext(c, requestMetadata)
	} else {
		md := metadata.Pairs(key, value)
		c = metadata.NewOutgoingContext(c, md)
	}
	return c
}

func GetSession(c context.Context) (sess *Session, exist bool) {
	if ctx, ok := GetContext(c); ok {
		if s := ctx.GetSession(); s != nil {
			sess = s
			exist = true
			return
		}
	}
	sess, exist = SessPool.GetUserSession(GetUID(c))
	return
}

func SetSessionValue(c context.Context, key string, value interface{}) (success bool) {
	if sess, ok := GetSession(c); ok {
		sess.SetSessValue(key, value)
		success = true
	}
	return
}

func GetSessionValue(c context.Context, key string) (ret interface{}, exist bool) {
	if sess, ok := GetSession(c); ok {
		if ret = sess.GetSessValue(key); ret != nil {
			exist = true
		}
	}
	return
}

func GetRequest(ctx context.Context) IProto {
	if ctx == nil {
		return nil
	}
	if iReq := ctx.Value(KeyNameRequest); iReq != nil {
		if req, ok := iReq.(IProto); ok {
			return req
		}
	}
	return nil
}

func GetResponse(ctx context.Context) IProto {
	if ctx == nil {
		return nil
	}
	if iRsp := ctx.Value(KeyNameResponse); iRsp != nil {
		if rsp, ok := iRsp.(IProto); ok {
			return rsp
		}
	}
	return nil
}

func GetError(ctx context.Context) error {
	if ctx == nil {
		return nil
	}
	if iError := ctx.Value(KeyNameError); iError != nil {
		if err, ok := iError.(error); ok {
			return err
		}
	}
	return nil
}

func GetBody(ctx context.Context) []byte {
	if ctx == nil {
		return nil
	}
	if iBody := ctx.Value(GinKeyBody); iBody != nil {
		if data, ok := iBody.([]byte); ok {
			return data
		}
	}
	return nil
}

func GetBeginTime(ctx context.Context) *time.Time {
	if ctx == nil {
		return nil
	}
	if iBegin := ctx.Value(ContextBeginKey); iBegin != nil {
		if begin, ok := iBegin.(*time.Time); ok {
			return begin
		}
	}
	return nil
}

func GetUID(ctx context.Context) string {
	if c, ok := GetContext(ctx); ok {
		return c.GetUID()
	}
	if uid, ok := GrpcGetUID(ctx); ok {
		return uid
	}
	if uid, ok := GinGetUID(ctx); ok {
		return uid
	}
	return ""
}

func GrpcGetUID(ctx context.Context) (uid string, bFound bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		us := md.Get(KeyNameUID)
		if len(us) < 1 {
			return
		}
		uid = us[0]
		bFound = true
		return
	}
	return
}

func GinGetUID(ctx context.Context) (uid string, ok bool) {
	var c *gin.Context
	if c, ok = GetGinContext(ctx); ok {
		var iUID interface{}
		if iUID, ok = c.Get(KeyNameUID); ok {
			if uid, ok = iUID.(string); ok {
				return
			}
		}
		if uid = c.Query(KeyNameUID); uid != "" {
			ok = true
		}
	}
	return
}

func GetBaseSpan(ctx context.Context) *tracetyp.Span {
	return tracetyp.GetBaseSpan(ctx)
}

func GetBaseCtx(ctx context.Context) context.Context {
	return tracetyp.GetBaseCtx(ctx)
}

// GetGinCurrentCtx 常在gin中间件调用，获取当前的ctx
func GetGinCurrentCtx(ctx context.Context) context.Context {
	iCtx := ctx.Value(GinKeyCurrentCtx)
	if iCtx != nil {
		if c, ok := iCtx.(context.Context); ok {
			return c
		}
	}
	return ctx
}

func GetRpcFullName(ctx context.Context) string {
	if gic, ok := GetGinContext(ctx); ok {
		return gic.FullPath()
	}
	if im := ctx.Value(KeyNameFullMethod); im != nil {
		return im.(string)
	}
	if fc, ok := GetContext(ctx); ok {
		return strconv.Itoa(int(fc.Cmd))
	}
	return ""
}

func GetRemoteIP(ctx context.Context) (addr string) {
	addr = strings.Split(GetRemoteAddr(ctx), ":")[0]
	return
}

func GetRemoteSrv(ctx context.Context) (srv string) {
	var ok bool
	if iSrv := ctx.Value(KeyNameRemoteSrv); iSrv != nil {
		if srv, ok = iSrv.(string); ok {
			return
		}
	}
	//读Meta data
	if srv, ok = GetMetaValue(ctx, KeyNameRemoteSrv); ok {
		return
	}
	return
}

func GetRequestID(ctx context.Context) string {
	if fc, ok := GetContext(ctx); ok {
		if fc.ReqId != "" {
			return fc.ReqId
		}
	}
	if val, ok := GetMetaValue(ctx, KeyNameReqID); ok {
		return val
	}
	return ""
}

func GetService(ctx context.Context) (srv interface{}, ok bool) {
	if srv = ctx.Value(KeyNameService); srv != nil {
		ok = true
	}
	return
}

func SetReqID(ctx context.Context, reqID string) context.Context {
	if fc, ok := GetContext(ctx); ok {
		fc.ReqId = reqID
	}
	ctx = SetMetaValue(ctx, KeyNameReqID, reqID)
	return ctx
}

func SetDiscoveryChanged(cb ftyps.TDiscoveryChanged) {
	fvars.OnDiscoveryCall = append(fvars.OnDiscoveryCall, cb)
}

func IsIProtoValid(pro IProto) bool {
	if pro == nil {
		return false
	}
	v := reflect.ValueOf(pro)
	if v.IsNil() {
		return false
	}
	return true
}

func IsLocal(ctx context.Context) bool {
	if isLocal, ok := ctx.Value(basedef.KeyNameIsLocal).(string); ok {
		return isLocal == "local"
	}
	return false
}

func GetRemoteAddr(ctx context.Context) (addr string) {
	//grpc
	p, ok := peer.FromContext(ctx)
	if ok {
		addr = p.Addr.String()
		return
	}
	//gin
	ginc, ok := GetGinContext(ctx)
	if ok {
		addr = util.GetAddrFromRequest(ginc.Request)
		return
	}
	//frpc
	if fc, ok := GetContext(ctx); ok {
		return fc.fromAddr
	}
	return
}

func GetRpcName(ctx context.Context) string {
	rpcName, _ := ctx.Value(basedef.KeyNameRpcName).(string)
	return rpcName
}
