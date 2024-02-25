package tracetyp

import (
	"bytes"
	"github.com/zhanmengao/gf/errors"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"sync"
	"sync/atomic"
)

// Span 用于兼容性和扩展
type Span struct {
	lock sync.RWMutex
	trace.Span
	attrList        []interface{} //当前Span的属性列表
	childSpan       []*Span       //记录所有孩子节点。非root节点的列表为空
	isRoot          bool          //是否是根节点
	isReport        bool          //表示当前span是否上报。当父span中有需要上报的child span，则会全部上报
	isAlreadyReport int32         //是否已经上报了
	root            *Span
}

// NewF9Span creates a span using default tracer.
func NewF9Span(span trace.Span, rootSpan *Span, isReport bool) *Span {
	sp := &Span{
		Span:      span,
		isRoot:    rootSpan == nil,
		isReport:  isReport,
		attrList:  make([]interface{}, 0, 10),
		childSpan: make([]*Span, 0, 10),
		root:      rootSpan,
	}
	rootSpan.AppendChild(sp)
	return sp
}

func (info *Span) AddAttr(key, value string) {
	if info.IsInvalid() {
		return
	}
	if value == "" {
		return
	}
	info.SetAttributes(key, value)
}

func (info *Span) SetAttributes(k string, v interface{}) {
	if info.IsInvalid() {
		return
	}
	info.lock.Lock()
	defer info.lock.Unlock()

	info.attrList = append(info.attrList, k)
	info.attrList = append(info.attrList, v)
}

func (info *Span) SpanString() string {
	if info.IsInvalid() {
		return ""
	}
	info.lock.RLock()
	defer info.lock.RUnlock()

	sb := bytes.Buffer{}
	for i := 0; i < len(info.attrList); i += 2 {
		if k, ok := info.attrList[i].(string); ok {
			if i+1 < len(info.attrList) {
				sb.WriteString(k)
				sb.WriteString(":")
				sb.WriteString(getInterface(info.attrList[i+1]))
				sb.WriteString(";")
			}
		}
	}
	return sb.String()
}

func (info *Span) SetError(err error) {
	if info.IsInvalid() {
		return
	}
	if err == nil {
		return
	}
	//error情况必定上报
	info.SetReport(true)
	//设置父span为true
	if info.root != nil {
		info.root.SetReport(true)
	}
	info.Span.SetStatus(codes.Error, err.Error())
	info.Span.RecordError(err)
	info.SetAttributes(TraceErrorCode, errors.Code(err))
	info.SetAttributes(TraceErrorMsg, errors.Message(err))
}

func (info *Span) End() {
	if info.IsInvalid() {
		return
	}
	info.lock.Lock()
	defer info.lock.Unlock()

	//如果是父节点已经报了，那么后续的节点也要上报
	if !info.isRoot && atomic.LoadInt32(&info.root.isAlreadyReport) >= 1 {
		info.end()
	}
	//如果父span确定会报，那么这里也提前结束子span
	if !info.isRoot && info.root.isReport {
		info.end()
	}
	//只有父节点可以调用End
	if !info.isRoot {
		return
	}
	//看看孩子是否有Report状态的
	for _, child := range info.childSpan {
		if child.IsReport() {
			info.isReport = true
			break
		}
	}
	count := 0
	//上报自己和所有孩子
	if info.isReport {
		info.end()
		for _, child := range info.childSpan {
			if count < 200 {
				child.end()
				count++
			} else {
				break
			}
		}
	}
}

func (info *Span) end() {
	if info.IsInvalid() {
		return
	}
	if !atomic.CompareAndSwapInt32(&info.isAlreadyReport, 0, 1) {
		return
	}
	for i := 0; i < len(info.attrList); i += 2 {
		if k, ok := info.attrList[i].(string); ok {
			if i+1 < len(info.attrList) {
				info.Span.SetAttributes(GetAttributes(k, info.attrList[i+1]))
			}
		}
	}
	info.Span.End()
}

func (info *Span) IsValid() bool {
	if info == nil || info.Span == nil {
		return false
	}
	//过滤掉non trace
	if _, ok := info.Span.(sdktrace.ReadOnlySpan); !ok {
		return false
	}
	return true
}

func (info *Span) IsInvalid() bool {
	return !info.IsValid()
}

func (info *Span) ToRpcSpan() *RpcSpan {
	return &RpcSpan{
		span: info,
	}
}

func (info *Span) ToMysqlSpan() *MysqlSpan {
	return &MysqlSpan{
		span: info,
	}
}

func (info *Span) ToRedisSpan() *RedisSpan {
	return &RedisSpan{
		span: info,
	}
}

func (info *Span) AppendChild(sp *Span) {
	if info.IsInvalid() {
		return
	}
	info.lock.Lock()
	defer info.lock.Unlock()
	info.childSpan = append(info.childSpan, sp)
}

func (info *Span) IsReport() bool {
	if info.IsInvalid() {
		return false
	}
	info.lock.RLock()
	defer info.lock.RUnlock()
	return info.isReport
}

func (info *Span) SetReport(report bool) {
	if info.IsInvalid() {
		return
	}
	info.lock.Lock()
	defer info.lock.Unlock()
	info.isReport = report
}
