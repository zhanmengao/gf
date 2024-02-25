package tracetyp

import (
	"bytes"
	"fmt"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"strings"
)

type TableSpan struct {
	span  *Span
	sheet string
}

func NewTableSpan(span *Span, sheet string) *TableSpan {
	return &TableSpan{
		span:  span,
		sheet: sheet,
	}
}

func (t *TableSpan) End(tabDefine, ab string, key, subKey interface{}, cfgList ...fmt.Stringer) {
	if t == nil || t.IsInvalid() {
		return
	}
	rs, ok := t.span.Span.(sdktrace.ReadOnlySpan)
	if !ok {
		return
	}
	buf := bytes.Buffer{}
	if ab != "" {
		buf.WriteString(fmt.Sprintf(TraceTableAB, ab))
	}
	if key != "" {
		buf.WriteString(fmt.Sprintf(TraceTableKey, key))
	}
	if subKey != "" {
		buf.WriteString(fmt.Sprintf(TraceTableSubKey, subKey))
	}
	//去重
	prefix := fmt.Sprintf(TraceReadTablePrefix, t.sheet)
	value := buf.String()

	for _, attr := range rs.Attributes() {
		if strings.HasPrefix(string(attr.Key), prefix) {
			if value == attr.Value.AsString() {
				return
			}
		}
	}
	t.span.AddAttr(fmt.Sprintf(TraceReadTable, prefix), value)
}

func (t *TableSpan) IsValid() bool {
	return t != nil && t.span.IsValid()
}

func (t *TableSpan) IsInvalid() bool {
	return !t.IsValid()
}
