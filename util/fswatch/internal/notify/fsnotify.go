package notify

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/zhanmengao/gf/errors"
	"github.com/zhanmengao/gf/util"
	"github.com/zhanmengao/gf/util/fswatch/internal/fsutil"
	"github.com/zhanmengao/gf/util/fswatch/wopt"
	"log/slog"
	"os"
	"path"
	"sync"
)

type FSNotify struct {
	watch    *fsnotify.Watcher
	listener map[string][]wopt.TEventCall
	lock     sync.RWMutex
	closeCh  chan struct{}
	once     sync.Once
	ctx      context.Context
}

func NewFSNotify(ctx context.Context) (f *FSNotify, err error) {
	f = &FSNotify{
		closeCh:  make(chan struct{}, 1),
		listener: make(map[string][]wopt.TEventCall, 5),
		ctx:      ctx,
	}
	if f.watch, err = fsnotify.NewWatcher(); err != nil {
		return
	}
	return
}

func (fn *FSNotify) AddWatcher(filePath string, cb wopt.TEventCall, autoCreate bool) (err error) {
	var ll []wopt.TEventCall
	fn.lock.Lock()
	defer fn.lock.Unlock()
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

	if ll, exist = fn.listener[filePath]; exist {
		fn.listener[filePath] = append(fn.listener[filePath], cb)
	} else {
		//先尝试add该文件
		if err = fn.watch.Add(filePath); err != nil {
			return
		}
		ll = append(ll, cb)
		fn.listener[filePath] = ll
	}
	return
}

func (fn *FSNotify) Run() {
	fn.once.Do(func() {
		go func() {
			defer fn.clear()
			for {
				select {
				case e := <-fn.watch.Events:
					fn.onEvent(e)
				case e := <-fn.watch.Errors:
					slog.ErrorContext(fn.ctx, "fsnotify err = %s ", e.Error())
				case <-fn.closeCh:
					return
				}
			}
		}()
	})
}

func (fn *FSNotify) onEvent(e fsnotify.Event) {
	dir, name := path.Split(e.Name)
	ec := wopt.FileWatchEvent{
		Dir:      dir,
		FileName: name,
	}
	if len(dir) > 0 {
		if dir[len(dir)-1] == os.PathSeparator {
			dir = dir[0 : len(dir)-1]
		}
	}
	dirEvent := false
	switch e.Op {
	case fsnotify.Write:
		ec.WatchType = wopt.WatchChanged
	case fsnotify.Remove:
		ec.WatchType = wopt.WatchDelete
		dirEvent = true
	case fsnotify.Create:
		ec.WatchType = wopt.WatchCreated
		dirEvent = true
	default:
		return
	}
	fn.lock.RLock()
	if ll, exist := fn.listener[e.Name]; exist {
		fn.lock.RUnlock()
		for _, v := range ll {
			util.SafeFunc(fn.ctx, func(ctx context.Context) {
				v(ec)
			})
		}
	} else {
		fn.lock.RUnlock()
	}
	if dirEvent {
		fn.lock.RLock()
		if ll, exist := fn.listener[dir]; exist {
			fn.lock.RUnlock()
			for _, v := range ll {
				util.SafeFunc(fn.ctx, func(ctx context.Context) {
					v(ec)
				})
			}
		} else {
			fn.lock.RUnlock()
		}
	}
}

func (fn *FSNotify) Close() {
	if fn.closeCh != nil {
		fn.closeCh <- struct{}{}
	}
}
func (fn *FSNotify) clear() {
	if fn.closeCh != nil {
		close(fn.closeCh)
		fn.closeCh = nil
	}
	if fn.watch != nil {
		fn.watch.Close()
		fn.watch = nil
	}
}
