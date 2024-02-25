package util

import (
	"context"
	"github.com/huandu/go-clone"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"github.com/zhanmengao/gf/util"
	"reflect"
)

type listener struct {
	key          string
	config       interface{}
	listenerList []cfgtp.WatchConfigCall
}

func newListener(key string, cfg interface{}) *listener {
	l := &listener{
		config:       cfg,
		listenerList: make([]cfgtp.WatchConfigCall, 0, 5),
		key:          key,
	}
	return l
}

func (l *listener) AppendListener(cb cfgtp.WatchConfigCall) {
	l.listenerList = append(l.listenerList, cb)
}

func (l *listener) Notify(ctx context.Context, cfg interface{}) {
	//一样，不返回
	if reflect.DeepEqual(l.config, cfg) {
		return
	}
	old := l.config
	for _, v := range l.listenerList {
		util.SafeFunc(ctx, func(ctx context.Context) {
			v(l.key, old, cfg)
		})
	}
	l.config = clone.Clone(cfg)
}
