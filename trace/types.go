package trace

import (
	"context"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
)

type fTraceLog struct {
}

func (f fTraceLog) Error(msg string) {
	slog.Error("[tracer] : ", msg)
}

func (f fTraceLog) Infof(msg string, args ...interface{}) {
	slog.Info("[tracer] : ", msg, args)
}

type spanExporter struct {
	sdktrace.SpanExporter
}

func (b *spanExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, s := range spans {
		// 监控 必须看那么是否唯一
		_ = s
		promCounterVec.WithLabelValues().Inc()
	}
	return b.SpanExporter.ExportSpans(ctx, spans)
}

func (b *spanExporter) Shutdown(ctx context.Context) error {
	return b.SpanExporter.Shutdown(ctx)
}
