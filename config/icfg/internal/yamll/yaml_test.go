package yamll

import (
	"github.com/zhanmengao/gf/config/icfg/cfgtp"
	"log"
	"testing"
	"time"
)

func TestYaml(t *testing.T) {
	cfg, err := NewConfigYaml("midplat.yaml")
	if err != nil {
		t.Fatal(err)
	}
	baseCfg := &cfgtp.BasicConfig{}
	if err = cfg.Get("basecfg", baseCfg); err != nil {
		t.Fatal(err)
	}
	t.Log(baseCfg)
}

func TestConfigYaml_GetAndWatch(t *testing.T) {
	cfg, err := NewConfigYaml("midplat.yaml")
	if err != nil {
		t.Fatal(err)
	}
	baseCfg := &cfgtp.BasicConfig{}
	if err = cfg.GetAndWatch("basecfg", baseCfg, func(key string, old, newConfig interface{}) {
		log.Printf("old = %v .new = %v ", old, newConfig)
	}); err != nil {
		t.Fatal(err)
	}
	t.Log(baseCfg)
	time.Sleep(time.Duration(1) * time.Minute)
}
