package yamll

import (
	"context"
	"github.com/huandu/go-clone"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"github.com/zhanmengao/gf/config/icfg/internal/util"
	"github.com/zhanmengao/gf/errors"
	"github.com/zhanmengao/gf/util/fswatch"
	"github.com/zhanmengao/gf/util/fswatch/wopt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log/slog"
	"reflect"
)

type ConfigYaml struct {
	filePath string
	watch    fswatch.IWatcher
	listener *util.ListenerKeys
	ctx      context.Context
}

// NewConfigYaml filePath表示yaml文件
func NewConfigYaml(ctx context.Context, filePath string) (cfg *ConfigYaml, err error) {
	cfg = &ConfigYaml{
		filePath: filePath,
		listener: util.NewListenerKeys(),
		ctx:      ctx,
	}
	cfg.watch = fswatch.DefaultWatcher
	if err = cfg.watch.AddWatcher(cfg.filePath, cfg.onFileChanged, false); err != nil {
		return
	}
	cfg.watch.Run()
	return
}

func (c *ConfigYaml) Get(key string, cfg interface{}) (err error) {
	body, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		return
	}
	kv := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(body, &kv); err != nil {
		return
	}
	val, exist := kv[key]
	if !exist {
		err = errors.NotFound("yaml key ", key)
		return
	}
	if body, err = yaml.Marshal(val); err != nil {
		return
	}
	if err = yaml.Unmarshal(body, cfg); err != nil {
		return
	}
	return
}
func (c *ConfigYaml) GetAndWatch(key string, cfg interface{}, cb cfgtp.WatchConfigCall) (err error) {
	if err = c.Get(key, cfg); err != nil {
		return
	}
	newConfig := clone.Clone(cfg)
	c.listener.AddListener(key, newConfig, cb)
	return
}

func (c *ConfigYaml) onFileChanged(e wopt.FileWatchEvent) {
	body, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		return
	}
	kv := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(body, &kv); err != nil {
		return
	}
	for key, val := range kv {
		name := key.(string)
		if cfg := c.listener.GetConfig(name); cfg != nil {
			cfg = reflect.New(reflect.TypeOf(cfg).Elem()).Interface()
			if body, err = yaml.Marshal(val); err != nil {
				slog.ErrorContext(c.ctx, "yaml marshal %v . %v error = %s ", name, val, err.Error())
				continue
			}
			if err = yaml.Unmarshal(body, cfg); err != nil {
				slog.ErrorContext(c.ctx, "yaml Unmarshal %v . %v error = %s ", name, val, err.Error())
				continue
			}
			c.listener.Notify(c.ctx, name, cfg)
		}
	}
}
