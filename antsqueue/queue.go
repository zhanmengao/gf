package antsqueue

import (
	"context"
	"github.com/zhanmengao/gf/errors"
	"github.com/zhanmengao/gf/util"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"
)

type AntsJobQueue struct {
	queue  *ants.PoolWithFunc
	job    chan *UserRequestList //放user的channel
	close  chan struct{}
	wg     sync.WaitGroup
	opt    *Option
	isOpen bool
	//user map
	userMap sync.Map
	tm      *time.Timer
	ctx     context.Context
}

func NewAntsJobQueue(opt ...*Option) *AntsJobQueue {
	ret := &AntsJobQueue{
		job:    make(chan *UserRequestList),
		close:  make(chan struct{}, 0),
		isOpen: true,
	}
	if len(opt) > 0 {
		ret.opt = opt[0]
	} else {
		ret.opt = NewOption()
	}
	opts := make([]ants.Option, 0)
	opts = append(opts, ants.WithPanicHandler(ret.panicHandler))
	if ret.opt.NonBlock {
		opts = append(opts, ants.WithNonblocking(true))
	}
	ret.queue, _ = ants.NewPoolWithFunc(ret.opt.Size, ret.do, opts...)
	util.Go(context.Background(), ret.run)
	return ret
}

func (q *AntsJobQueue) do(iUser interface{}) {
	//实际执行任务的地方
	u := iUser.(*UserRequestList)
	for {
		select {
		case f := <-u.request:
			f.Func(f.Ctx)
		default:
			return
		}
	}
}

func (q *AntsJobQueue) run(ctx context.Context) {
	q.isOpen = true
	defer q.queue.Release()
	q.tm = time.NewTimer(time.Duration(5) * time.Second)
	for {
		select {
		case u := <-q.job:
			//注意这里取的是user，也就是说只要阻塞的用户不多，是不会卡到其他用户的
			if err := q.queue.Invoke(u); err != nil {
				slog.ErrorContext(ctx, "u[%s] busy", u.key)
			}
		case <-q.close:
			//取出所有消息，处理完return
			q.clear()
			q.wg.Wait()
		case <-q.tm.C:
			q.tm = time.NewTimer(time.Duration(5) * time.Second)
		}
	}
}

func (q *AntsJobQueue) PushJob(ctx context.Context, key string, f func(ctx2 context.Context)) (err error) {
	if !q.isOpen {
		err = errors.BadRequest("", "queue stop")
		return
	}
	u := NewUserRequestList(key, 20)
	if iu, exist := q.userMap.LoadOrStore(key, u); exist {
		u = iu.(*UserRequestList)
	} else {
		u.Init()
	}
	err = q.pushJob(u, &TJob{
		Func: f,
		Ctx:  ctx,
	})
	return
}

func (q *AntsJobQueue) pushJob(u *UserRequestList, f *TJob) (err error) {
	//先放进该用户的请求列表
	if q.opt.NonBlock {
		select {
		case u.request <- f:
		default:
			err = errors.BadRequest("", "busy")
			return
		}
	} else {
		u.request <- f
	}

	//把该用户放进队列的请求列表
	if q.opt.NonBlock {
		select {
		case q.job <- u:
		default:
			err = errors.BadRequest("", "busy ")
			return
		}
	} else {
		q.job <- u
	}

	return
}

func (q *AntsJobQueue) clearUser() {
	now := time.Now().Unix()
	q.userMap.Range(func(iUID interface{}, iUser interface{}) bool {
		user := iUser.(*UserRequestList)
		if now-user.lastActive > 5 {
			q.userMap.Delete(iUID)
		}
		return true
	})
}

func (q *AntsJobQueue) Close() {
	q.isOpen = false
	if q.close != nil {
		q.close <- struct{}{}
	}
}

func (q *AntsJobQueue) clear() {
	for {
		select {
		case u := <-q.job:
			//如果是非阻塞模式，这里会返回err
			if err := q.queue.Invoke(u); err != nil {
				slog.Error("clear error %s", err.Error())
			}
		default:
			close(q.job)
			return
		}
	}
}

func (q *AntsJobQueue) panicHandler(iu interface{}) {
	slog.ErrorContext(q.ctx, "panic", string(debug.Stack()))
}
