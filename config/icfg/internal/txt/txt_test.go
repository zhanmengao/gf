package txt

import (
	"context"
	"log"
	"testing"
)

var c *ConfigTxt
var ctx = context.Background()

func TestMain(m *testing.M) {
	var err error
	c, err = NewConfigTxt(ctx, ".")
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestGet(t *testing.T) {
	var txt []byte
	err := c.Get("test.txt", &txt)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(txt))
}

func TestGetAndWatch(t *testing.T) {
	var txt []byte
	err := c.GetAndWatch("test.txt", &txt, func(key string, oldConfig interface{}, newConfig interface{}) {
		log.Println(string(oldConfig.([]byte)), " ", string(newConfig.([]byte)))
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(txt))
	select {}
}
