package tracetyp

type MysqlSpan struct {
	span *Span
}

func NewMysqlSpan(span *Span) *MysqlSpan {
	return &MysqlSpan{
		span: span,
	}
}

func (t *MysqlSpan) End(tabName, sql string, rows int64, result interface{}, err error) {
	if t == nil || t.IsInvalid() {
		return
	}
	t.span.SetError(err)
	if result != nil {
		t.span.SetAttributes(TraceResult, result)
	}
	t.span.SetAttributes(TraceMysqlSql, sql)
	t.span.SetAttributes(TraceMysqlTable, tabName)
	t.span.SetAttributes(TraceMysqlRows, rows)
	t.span.End()
}
func (t *MysqlSpan) IsValid() bool {
	return t != nil && t.span.IsValid()
}

func (t *MysqlSpan) IsInvalid() bool {
	return !t.IsValid()
}
func (t *MysqlSpan) SetAttributes(k string, v interface{}) {
	if t.IsInvalid() {
		return
	}
	t.span.SetAttributes(k, v)
}

func (t *MysqlSpan) AddAttr(key, value string) {
	if t.IsInvalid() {
		return
	}
	if value == "" {
		return
	}
	t.SetAttributes(key, value)
}
