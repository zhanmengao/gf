package fstat

import (
	"context"
	"github.com/zhanmengao/gf/errors"
	"github.com/zhanmengao/gf/util"

	"github.com/zhanmengao/gf/util/fswatch/internal/fsutil"
	"github.com/zhanmengao/gf/util/fswatch/wopt"
	"os"
	"path"
	"sync"
	"time"
)

type FStat struct {
	fileMap map[string]*FileInfo
	lock    sync.RWMutex
	timer   *time.Timer
	closeCh chan struct{}
	once    sync.Once
	ctx     context.Context
}

func NewFStat(ctx context.Context) *FStat {
	return &FStat{
		fileMap: make(map[string]*FileInfo),
		closeCh: make(chan struct{}, 1),
		ctx:     ctx,
	}
}

func (f *FStat) AddWatcher(filePath string, cb wopt.TEventCall, autoCreate bool) (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	//检测文件是否存在
	var exist bool
	if exist, err = fsutil.Exists(filePath); err != nil {
		return
	} else if !exist {
		if autoCreate {
			if err = fsutil.TouchFile(filePath); err != nil {
				return
			}
		} else {
			err = errors.NotFound("file not found ", filePath).SetBasicErr(err)
			return
		}
	}

	if fi, exist := f.fileMap[filePath]; exist {
		fi.callList = append(fi.callList, cb)
	} else {
		dir, name := path.Split(filePath)
		var info os.FileInfo
		//stat获取不到也加入listen，如果下次获取到了则触发create事件
		if info, err = os.Stat(filePath); err != nil {
			err = nil
		}
		f.fileMap[filePath] = &FileInfo{
			callList: make([]wopt.TEventCall, 0),
			dir:      dir,
			stat:     info,
			name:     name,
		}
		f.fileMap[filePath].callList = append(f.fileMap[filePath].callList, cb)
	}
	return
}

func (f *FStat) Run() {
	f.once.Do(func() {
		f.timer = time.NewTimer(time.Duration(1) * time.Second)
		go func() {
			for {
				select {
				case <-f.timer.C:
					f.tick()
					f.timer.Reset(time.Duration(1) * time.Second)
				case <-f.closeCh:
					return
				}
			}
		}()
	})
}

func (f *FStat) tick() {
	f.lock.RLock()
	defer f.lock.RUnlock()
	for name, v := range f.fileMap {
		stat, _ := os.Stat(name)
		mt := wopt.WatchInit
		if stat == nil && v.stat == nil {
			continue
		} else if stat == nil && v.stat != nil {
			mt = wopt.WatchDelete
		} else if stat != nil && v.stat == nil {
			mt = wopt.WatchCreated
		} else if stat.ModTime() != v.stat.ModTime() {
			mt = wopt.WatchChanged
		}
		v.stat = stat
		ev := wopt.FileWatchEvent{
			WatchType: mt,
			FileName:  v.name,
			Dir:       v.dir,
		}
		if mt != wopt.WatchInit {
			for _, fun := range v.callList {
				util.SafeFunc(f.ctx, func(ctx context.Context) {
					fun(ev)
				})
			}
		}
	}
}

func (f *FStat) Close() {
	if f.closeCh != nil {
		f.closeCh <- struct{}{}
	}
}

func (f *FStat) clear() {
	if f.closeCh != nil {
		close(f.closeCh)
		f.closeCh = nil
	}
}
