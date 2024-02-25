package fstat

import (
	"context"
	"github.com/zhanmengao/gf/util/fswatch/wopt"
	"log"
	"testing"
)

var notify *FStat
var ctx = context.Background()

func TestMain(m *testing.M) {
	var err error
	notify = NewFStat(ctx)
	if err != nil {
		panic(err)
	}
	defer notify.Close()
	m.Run()
}
func eventCB(e wopt.FileWatchEvent) {
	log.Println(e)
}
func TestAddNotify(t *testing.T) {
	err := notify.AddWatcher("proc.go", eventCB)
	if err != nil {
		t.Fatal(err)
	}
	if err = notify.AddWatcher("temp.txt", eventCB); err != nil {
		t.Fatal(err)
	}
	notify.Run()
	notify.Close()
}

func TestWatch(t *testing.T) {
	err := notify.AddWatcher("proc.go", eventCB)
	if err != nil {
		t.Fatal(err)
	}
	if err = notify.AddWatcher("temp.txt", eventCB); err != nil {
		t.Fatal(err)
	}
	notify.Run()
	select {}
}
