package util

import (
	"context"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"sync"
)

type ListenerKeys struct {
	listenerMap sync.Map
}

func NewListenerKeys() *ListenerKeys {
	return &ListenerKeys{}
}

func (lk *ListenerKeys) AddListener(key string, cfg interface{}, cb cfgtp.WatchConfigCall) {
	l := newListener(key, cfg)
	if iL, exist := lk.listenerMap.LoadOrStore(key, l); exist {
		l = iL.(*listener)
	}
	l.AppendListener(cb)
}

func (lk *ListenerKeys) Notify(ctx context.Context, key string, cfg interface{}) (exist bool) {
	var iL interface{}
	if iL, exist = lk.listenerMap.Load(key); exist {
		l := iL.(*listener)
		l.Notify(ctx, cfg)
	}
	return
}

func (lk *ListenerKeys) GetConfig(key string) (cfg interface{}) {
	if iL, exist := lk.listenerMap.Load(key); exist {
		l := iL.(*listener)
		cfg = l.config
	}
	return
}
