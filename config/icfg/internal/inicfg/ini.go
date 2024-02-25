package inicfg

import (
	"context"
	"github.com/huandu/go-clone"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"github.com/zhanmengao/gf/config/icfg/internal/util"
	"github.com/zhanmengao/gf/util/fswatch"
	"github.com/zhanmengao/gf/util/fswatch/wopt"
	"gopkg.in/ini.v1"
	"log/slog"
	"reflect"
)

type ConfigINI struct {
	filePath string
	file     *ini.File
	listener *util.ListenerKeys
	watch    fswatch.IWatcher
	ctx      context.Context
}

// NewConfigINI filePath表示配置文件
func NewConfigINI(ctx context.Context, filePath string) (cfg *ConfigINI, err error) {
	cfg = &ConfigINI{
		filePath: filePath,
		listener: util.NewListenerKeys(),
		ctx:      ctx,
	}
	if cfg.file, err = ini.Load(cfg.filePath); err != nil {
		return
	}
	cfg.watch = fswatch.DefaultWatcher
	if err = cfg.watch.AddWatcher(cfg.filePath, cfg.onFileChanged, false); err != nil {
		return
	}
	cfg.watch.Run()
	return
}

func (c *ConfigINI) Get(key string, cfg interface{}) (err error) {
	var s *ini.Section
	if s, err = c.file.GetSection(key); err != nil {
		return
	}
	keys := s.Keys()
	rv := reflect.ValueOf(cfg).Elem()
	//给每个字段赋值
	for _, k := range keys {
		filter := rv.FieldByName(k.Name())
		if !filter.IsValid() {
			continue
		}
		if err = util.SetField(filter, k.Value()); err != nil {
			return
		}
	}
	return
}
func (c *ConfigINI) GetAndWatch(key string, cfg interface{}, cb cfgtp.WatchConfigCall) (err error) {
	if err = c.Get(key, cfg); err != nil {
		return
	}
	newConfig := clone.Clone(cfg)
	c.listener.AddListener(key, newConfig, cb)
	return
}

func (c *ConfigINI) onFileChanged(e wopt.FileWatchEvent) {
	if err := c.file.Reload(); err != nil {
		slog.ErrorContext(c.ctx, "load config : %s . error = %s ", c.filePath, err.Error())
		return
	}
	sl := c.file.Sections()
	for _, k := range sl {
		name := k.Name()
		if name == "DEFAULT" {
			name = ""
		}
		cfg := c.listener.GetConfig(name)
		if cfg != nil {
			cfg = reflect.New(reflect.TypeOf(cfg).Elem()).Interface()
			if err := c.Get(name, cfg); err != nil {
				slog.ErrorContext(c.ctx, "get error = %s ", err.Error())
			} else {
				c.listener.Notify(c.ctx, name, cfg)
			}
		}
	}
	return
}
