package vipercfg

import (
	"github.com/fsnotify/fsnotify"
	"github.com/huandu/go-clone"
	"github.com/spf13/viper"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"reflect"
	"sync"
)

type ViperCfg struct {
	Viper        *viper.Viper
	OldValue     sync.Map
	CallbackList map[string][]cfgtp.WatchConfigCall
	WatchStatus  bool
}

func NewViperCfg(cfgfile string, filetype string) (*ViperCfg, error) {
	v := viper.New()
	v.SetConfigFile(cfgfile)
	v.SetConfigType(filetype)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &ViperCfg{
		Viper:        v,
		CallbackList: make(map[string][]cfgtp.WatchConfigCall),
	}, nil
}

func (y *ViperCfg) Get(key string, cfg interface{}) (err error) {
	return y.Viper.UnmarshalKey(key, cfg)
}

func (y *ViperCfg) GetAndWatch(key string, cfg interface{}, cb cfgtp.WatchConfigCall) (err error) {
	if err = y.Get(key, cfg); err != nil {
		return err
	}

	oldcfg := clone.Clone(cfg)
	y.AddWatcher(key, oldcfg, cb)

	y.Watch()
	return nil
}

func (y *ViperCfg) AddWatcher(key string, oldcfg interface{}, cb cfgtp.WatchConfigCall) {
	y.OldValue.Store(key, oldcfg)
	y.CallbackList[key] = append(y.CallbackList[key], cb)
}

func (y *ViperCfg) Watch() {
	if !y.WatchStatus {
		y.WatchStatus = true

		y.Viper.OnConfigChange(func(in fsnotify.Event) {
			if in.Op == fsnotify.Write {
				for k, cbs := range y.CallbackList {
					v, ok := y.OldValue.Load(k)
					if ok {
						T := reflect.TypeOf(v)
						newcfg := reflect.New(T.Elem()).Interface()
						y.Viper.UnmarshalKey(k, newcfg)
						if reflect.TypeOf(v) == reflect.TypeOf(newcfg) && !reflect.DeepEqual(v, newcfg) {
							for _, cb := range cbs {
								cb(k, v, newcfg)
							}
						}

						y.OldValue.Store(k, clone.Clone(newcfg))
					}
				}
			}
		})

		y.Viper.WatchConfig()
	}
}
