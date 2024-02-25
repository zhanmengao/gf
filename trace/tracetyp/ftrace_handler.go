package tracetyp

import (
	"bytes"
	"fmt"
)

type RpcSpan struct {
	span *Span
}

func NewRpcSpan(span *Span) *RpcSpan {
	return &RpcSpan{
		span: span,
	}
}

func (h *RpcSpan) End(req fmt.Stringer, rsp fmt.Stringer, err error) {
	if h == nil || h.IsInvalid() {
		return
	}
	//h.span.AddEvent("message", trace.WithAttributes(GetAttributes(TraceReq, req)))
	//h.span.AddEvent("message", trace.WithAttributes(GetAttributes(TraceResult, rsp)))
	h.span.SetAttributes(TraceReq, req)
	h.span.SetAttributes(TraceResult, rsp)
	h.span.SetError(err)
	h.span.End()
}

func (h *RpcSpan) SetRemoteAddr(addr string) {
	if h == nil || h.IsInvalid() {
		return
	}
	h.span.SetAttributes(TraceRemoteAddr, addr)
}

func (h *RpcSpan) SetRemoteSrv(srv string) {
	if h == nil || h.IsInvalid() {
		return
	}
	h.span.SetAttributes(TraceRemoteSrv, srv)
}

func (h *RpcSpan) SetUID(uid string) {
	if h == nil || h.IsInvalid() {
		return
	}
	h.span.SetAttributes(TraceUID, uid)
}

func (h *RpcSpan) SetIsLocal(isLocal bool) {
	if h == nil || h.IsInvalid() {
		return
	}
	h.span.SetAttributes(TraceRpcIsLocal, isLocal)
}

func (h *RpcSpan) SetBodySize(sz int) {
	if h == nil || h.IsInvalid() {
		return
	}
	h.span.SetAttributes(TraceBodySize, sz)
}

func (h *RpcSpan) SetHeader(header map[string]string) {
	if h == nil || h.IsInvalid() {
		return
	}
	headStr := bytes.Buffer{}
	for k, v := range header {
		headStr.WriteString(k)
		headStr.WriteByte(':')
		headStr.WriteString(v)
		headStr.WriteByte('\t')
	}
	h.span.AddAttr(TraceHttpHeader, headStr.String())
}
func (h *RpcSpan) IsValid() bool {
	return h != nil && h.span.IsValid()
}

func (h *RpcSpan) IsInvalid() bool {
	return !h.IsValid()
}

func (h *RpcSpan) SetAttributes(k string, v interface{}) {
	if h.IsInvalid() {
		return
	}
	h.span.SetAttributes(k, v)
}

func (h *RpcSpan) AddAttr(key, value string) {
	if h.IsInvalid() {
		return
	}
	if value == "" {
		return
	}
	h.SetAttributes(key, value)
}

func (h *RpcSpan) SetPlatform(platform string) {
	h.span.SetAttributes(TracePlatform, platform)
}

func (h *RpcSpan) SetChannel(channel string) {
	h.span.SetAttributes(TraceChannel, channel)
}

func (h *RpcSpan) SetSeqID(seq int32) {
	h.span.SetAttributes(TraceSeqID, seq)
}

func (h *RpcSpan) SetReqID(reqID string) {
	h.span.SetAttributes(TraceReqID, reqID)
}
