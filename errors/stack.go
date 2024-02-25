package errors

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

var (
	bufPool    sync.Pool
	stackBegin int = 3
	stackDepth int = 3
)

func init() {
	bufPool.New = func() interface{} {
		return &strings.Builder{}
	}
}

func stack(skip int) string {
	buf := bufPool.Get().(*strings.Builder)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	for i := stackBegin; i < stackDepth+stackBegin; i++ {
		pc, file, line, ok := runtime.Caller(i + skip)
		if ok {
			buf.WriteString(fmt.Sprintf("%s[%d]", file, line))
			buf.WriteRune(' ')
			buf.WriteString(runtime.FuncForPC(pc).Name())
			buf.WriteRune('\t')
		} else {
			break
		}
	}
	return buf.String()
}

func SetBeginStack(begin int) {
	stackBegin = begin
}

func SetMaxStack(depth int) {
	stackDepth = depth
}
