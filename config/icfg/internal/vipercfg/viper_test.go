package vipercfg

import (
	"fmt"
	"testing"
	"time"
)

type Server struct {
	Addr string
	Port int
}

func TestNewViperCfg(t *testing.T) {
	y, err := NewViperCfg("test.json", "json")
	if err != nil {
		t.Fatal(err)
	}

	c := Server{}

	y.GetAndWatch("Server", &c, func(key string, oldConfig interface{}, newConfig interface{}) {
		o, ok := oldConfig.(*Server)
		if ok {
			fmt.Println("Update Config Old Is:", o)
		}
		n, ok := newConfig.(*Server)
		if ok {
			fmt.Println("Update Config New Is:", n)
		}
	})

	fmt.Println("Get Config :", c)

	time.Sleep(time.Hour)

}
