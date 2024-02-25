package util

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strings"
	"time"
)

// Go runs a safe goroutine
func Go(ctx context.Context, f func(context.Context)) {
	if f == nil {
		return
	}
	go SafeFunc(ctx, f)
}

// SafeFunc safe function call
func SafeFunc(ctx context.Context, f func(c context.Context)) {
	defer func() {
		if r := recover(); r != nil {
			LogPanic(ctx, r)
		}
	}()
	f(ctx)
}

func LogPanic(ctx context.Context, r interface{}) {
	stack := string(debug.Stack())
	fmt.Println(time.Now().String())
	fmt.Println(r)
	fmt.Println(stack)
	slog.ErrorContext(ctx, "%v : %s ", r, strings.ReplaceAll(stack, "\n", "\t"))
}
