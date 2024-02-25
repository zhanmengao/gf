package inicfg

import (
	"context"
	"log"
	"testing"
)

var c *ConfigINI
var ctx = context.Background()

func TestMain(m *testing.M) {
	var err error
	c, err = NewConfigINI(ctx, "ops.ini")
	if err != nil {
		panic(err)
	}
	m.Run()
}

type tops struct {
	LogReportAddr   string
	InnerTimeout    int64
	CoordinatorAddr string
}

func TestINI(t *testing.T) {
	ops := &tops{}
	err := c.Get("", ops)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ops)
}

type FIMConfig struct {
	WsPort    int64
	ConnClose int64
}

func TestKey(t *testing.T) {
	fim := &FIMConfig{}
	err := c.Get("fimgate", fim)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fim)
}

func TestGetAndWatch(t *testing.T) {
	ops := &tops{}
	err := c.GetAndWatch("", ops, func(key string, oldConfig interface{}, newConfig interface{}) {
		log.Println(oldConfig, " ", newConfig)
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ops)
	select {}
}
