package fstat

import (
	"github.com/zhanmengao/gf/util/fswatch/wopt"
	"os"
)

type FileInfo struct {
	stat     os.FileInfo
	dir      string
	name     string
	callList []wopt.TEventCall
}
