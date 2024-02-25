package tracetyp

import (
	"fmt"
)

type RedisSpan struct {
	span *Span
}

func NewRedisSpan(span *Span) *RedisSpan {
	return &RedisSpan{
		span: span,
	}
}

func (t *RedisSpan) SetTTL(ttl int) {
	if t == nil || t.IsInvalid() {
		return
	}
	t.span.SetAttributes(TraceRedisTTL, ttl)
}

func (t *RedisSpan) End(keyFmt, key, fieldFmt, field string, exist, cache bool, result fmt.Stringer, err error) {
	if t == nil || t.IsInvalid() {
		return
	}
	t.span.AddAttr(TraceRedisKey, key)
	t.span.AddAttr(TraceRedisKeyFmt, keyFmt)
	if field != "" {
		t.span.AddAttr(TraceRedisField, field)
	}
	if fieldFmt != "" {
		t.span.AddAttr(TraceRedisFieldFmt, fieldFmt)
	}
	t.span.SetAttributes(TraceRedisExist, exist)
	t.span.SetAttributes(TraceRedisCache, cache)
	t.span.SetError(err)
	if result != nil {
		//t.span.AddEvent("message", trace.WithAttributes(GetAttributes(TraceResult, result)))
		t.span.SetAttributes(TraceResult, result)
	}
	t.span.End()
}
func (t *RedisSpan) IsValid() bool {
	return t != nil && t.span.IsValid()
}

func (t *RedisSpan) IsInvalid() bool {
	return !t.IsValid()
}

func (t *RedisSpan) SetAttributes(k string, v interface{}) {
	if t.IsInvalid() {
		return
	}
	t.span.SetAttributes(k, v)
}

func (t *RedisSpan) AddAttr(key, value string) {
	if t.IsInvalid() {
		return
	}
	if value == "" {
		return
	}
	t.SetAttributes(key, value)
}
