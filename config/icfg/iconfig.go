package icfg

import (
	"context"
	"fmt"
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"github.com/zhanmengao/gf/config/icfg/internal/inicfg"
	"github.com/zhanmengao/gf/config/icfg/internal/txt"
	"github.com/zhanmengao/gf/config/icfg/internal/vipercfg"
	"github.com/zhanmengao/gf/config/icfg/internal/yamll"
	"github.com/zhanmengao/gf/errors"
)

func NewConfig(ctx context.Context, filePath string, tp cfgtp.CfgType) (ret cfgtp.IConfig, err error) {
	switch tp {
	case cfgtp.INI:
		ret, err = inicfg.NewConfigINI(ctx, filePath)
	case cfgtp.TXT:
		ret, err = txt.NewConfigTxt(ctx, filePath)
	case cfgtp.Json:
		ret, err = vipercfg.NewViperCfg(filePath, "json")
	case cfgtp.YAML:
		ret, err = yamll.NewConfigYaml(ctx, filePath)
	case cfgtp.Toml:
		ret, err = vipercfg.NewViperCfg(filePath, "toml")
	case cfgtp.VYaml:
		ret, err = vipercfg.NewViperCfg(filePath, "yaml")
	default:
		err = errors.NotFound(fmt.Sprintf("config type %d ", tp), "")
	}
	return
}
