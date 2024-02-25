package fswatch

import (
	"context"
	"github.com/zhanmengao/gf/util/fswatch/internal/fstat"
	"github.com/zhanmengao/gf/util/fswatch/internal/notify"
	"github.com/zhanmengao/gf/util/fswatch/wopt"
)

type IWatcher interface {
	AddWatcher(filePath string, cb wopt.TEventCall, autoCreate bool) (err error) //加监控
	Run()
	Close()
}

var DefaultWatcher = NewWatcher(context.Background())

func NewWatcher(ctx context.Context) (w IWatcher) {
	var err error
	if w, err = notify.NewFSNotify(ctx); err != nil {
		w = fstat.NewFStat(ctx)
	}
	return
}
