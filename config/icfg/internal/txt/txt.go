package txt

import (
	"context"
	"github.com/huandu/go-clone"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"github.com/zhanmengao/gf/config/icfg/internal/util"
	"github.com/zhanmengao/gf/util/fswatch"
	"github.com/zhanmengao/gf/util/fswatch/wopt"
	"io/ioutil"
	"log/slog"
	"path"
)

type ConfigTxt struct {
	filePath string
	Content  []byte
	watch    fswatch.IWatcher
	listener *util.ListenerKeys
	ctx      context.Context
}

// NewConfigTxt filePath表示目录
func NewConfigTxt(ctx context.Context, filePath string) (cfg *ConfigTxt, err error) {
	cfg = &ConfigTxt{
		filePath: filePath,
		listener: util.NewListenerKeys(),
		ctx:      ctx,
	}
	cfg.watch = fswatch.DefaultWatcher
	if err = cfg.watch.AddWatcher(cfg.filePath, cfg.onFileChanged, true); err != nil {
		return
	}
	cfg.watch.Run()
	return
}

func (c *ConfigTxt) Get(key string, cfg interface{}) (err error) {
	filePath := path.Join(c.filePath, key)
	ret := cfg.(*[]byte)
	var content []byte
	if content, err = ioutil.ReadFile(filePath); err != nil {
		return
	}
	*ret = content
	return
}
func (c *ConfigTxt) GetAndWatch(key string, cfg interface{}, cb cfgtp.WatchConfigCall) (err error) {
	if err = c.Get(key, cfg); err != nil {
		return
	}
	data := *cfg.(*[]byte)
	newConfig := clone.Clone(data)
	if err = c.watch.AddWatcher(path.Join(c.filePath, key), c.onFileChanged, false); err != nil {
		return
	}
	c.listener.AddListener(key, newConfig, cb)
	return
}

func (c *ConfigTxt) onFileChanged(e wopt.FileWatchEvent) {
	var cfg []byte
	if err := c.Get(e.FileName, &cfg); err != nil {
		slog.ErrorContext(c.ctx, "read %s error = %s ", e.FileName, err.Error())
	}
	c.listener.Notify(c.ctx, e.FileName, cfg)
}
