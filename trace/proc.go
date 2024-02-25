package trace

import (
	"context"
	"fmt"
	"github.com/zhanmengao/gf/basedef"
	"github.com/zhanmengao/gf/trace/tracetyp"
	"go.opentelemetry.io/otel/trace"
)

// GetTracer 创建一个空tracer
func GetTracer(tracerName string) trace.Tracer {
	// Trace span start.
	tr := getTracerProvider(tracerName)
	return tr
}

func NewSpan(ctx context.Context, tracerName string, spanName string, opts ...trace.SpanStartOption) (context.Context, *tracetyp.Span) {
	return NewSpanByUse(ctx, tracerName, spanName, true, opts...)
}

func NewSpanByUse(ctx context.Context, tracerName string, spanName string, isReport bool, opts ...trace.SpanStartOption) (context.Context, *tracetyp.Span) {
	tr := GetTracer(tracerName)
	if tr == nil {
		return ctx, nil
	}
	//查看是否有父节点
	var rootSpan *tracetyp.Span
	if iRootSp := ctx.Value(basedef.KeyNameRootSpan); iRootSp != nil {
		rootSpan, _ = iRootSp.(*tracetyp.Span)
	}
	ctx, span := tr.Start(ctx, fmt.Sprintf("%s.%s", tracerName, spanName), opts...)
	f9sp := tracetyp.NewF9Span(span, rootSpan, isReport)
	//放到ctx里
	if rootSpan == nil {
		//root span
		ctx = context.WithValue(ctx, basedef.KeyNameRootSpan, f9sp)
		//current span
		ctx = context.WithValue(ctx, basedef.KeyNameCurrentSpan, f9sp)
		ctx = context.WithValue(ctx, basedef.KeyNameBaseSpanCtx, ctx)
	} else {
		//current span
		ctx = context.WithValue(ctx, basedef.KeyNameCurrentSpan, f9sp)
	}

	return ctx, f9sp
}

func NewTableSpan(ctx context.Context, sheetName string) *tracetyp.TableSpan {
	var span *tracetyp.Span
	//如果有Root Span，说明要进行tracer
	if root := tracetyp.GetBaseCtx(ctx); root != nil {
		span = tracetyp.GetCurrentSpan(ctx)
	}
	return tracetyp.NewTableSpan(span, sheetName)
}

func NewRedisSpan(ctx context.Context, dbName, command, keyDesc string) *tracetyp.RedisSpan {
	var span *tracetyp.Span
	//如果有Root Span，说明要进行tracer
	if root := tracetyp.GetBaseCtx(ctx); root != nil {
		_, span = NewSpanByUse(ctx, "redis", fmt.Sprintf("%s.%s.%s", dbName, command, keyDesc), false)
	}
	return tracetyp.NewRedisSpan(span)
}

// NewRpcSpan
//
//	@Description: 构造一个RPC Span
//	@param ctx 当前上下文
//	@param rpcName RPC名称，http/grpc/frpc
//	@param cmd 接口名
//	@param op 进行的操作，例如Request、Handler
//	@param isReport 请求全部填false，handler根据对端是否产生来生成
//	@param opts
//	@return context.Context
//	@return *tracetyp.RpcSpan
func NewRpcSpan(ctx context.Context, rpcName, cmd string, op tracetyp.RpcOperator, opts ...trace.SpanStartOption) (context.Context, *tracetyp.RpcSpan) {
	return NewRpcSpanByUse(ctx, rpcName, cmd, op, true, opts...)
}
func NewRpcSpanByUse(ctx context.Context, rpcName, cmd string, op tracetyp.RpcOperator, isReport bool, opts ...trace.SpanStartOption) (context.Context, *tracetyp.RpcSpan) {
	var span *tracetyp.Span
	switch op {
	case tracetyp.RpcOpHandler:
		opts = append(opts, trace.WithSpanKind(trace.SpanKindServer))
	case tracetyp.RpcOpSend, tracetyp.RpcOpRequest:
		opts = append(opts, trace.WithSpanKind(trace.SpanKindClient))
	case tracetyp.RpcMqConsumer:
		opts = append(opts, trace.WithSpanKind(trace.SpanKindConsumer))
	case tracetyp.RpcMqProducer:
		opts = append(opts, trace.WithSpanKind(trace.SpanKindProducer))
	}
	ctx, span = NewSpanByUse(ctx, fmt.Sprintf("%s.%s", rpcName, tracetyp.GetRpcOpName(op)), cmd, isReport, opts...)
	return ctx, tracetyp.NewRpcSpan(span)
}

func NewMysqlSpan(ctx context.Context, dbName, command, keyDesc string) (context.Context, *tracetyp.MysqlSpan) {
	var span *tracetyp.Span
	//如果有Root Span，说明要进行tracer
	if root := tracetyp.GetBaseCtx(ctx); root != nil {
		ctx, span = NewSpanByUse(ctx, "mysql", fmt.Sprintf("%s.%s.%s", dbName, command, keyDesc), false)
	}
	return ctx, tracetyp.NewMysqlSpan(span)
}
