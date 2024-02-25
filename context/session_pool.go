package framework

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"
)

const maxSessionExpire = 30 * 60 * 1000 // session在内存中最多存放30分钟

type SessionPool struct {
	sessionPool  sync.Map
	offlineHeart int64

	sessionExpireMillSecond int64
	ctx                     context.Context
}

func NewSessionPool(ctx context.Context) *SessionPool {
	sp := &SessionPool{
		offlineHeart:            60 * 1000,
		sessionExpireMillSecond: maxSessionExpire,
		ctx:                     ctx,
	}
	go func() {
		tk := time.NewTicker(1 * time.Minute)
		for range tk.C {
			sp.clearExpiredSession()
		}
	}()
	return sp
}

func (sp *SessionPool) loadSession(scope string, key string) (sess *Session, ok bool) {
	k := sessionKey(scope, key)
	data, ok := sp.sessionPool.Load(k)
	if !ok {
		return
	}
	sess = data.(*Session)
	return
}

func (sp *SessionPool) createSession(scope string, key string) (sess *Session, loaded bool) {
	s := newSession(&sp.offlineHeart, &sp.sessionExpireMillSecond, key)
	k := sessionKey(scope, key)
	act, loaded := sp.sessionPool.LoadOrStore(k, s)
	if loaded {
		sess = act.(*Session)
	} else {
		sess = s
	}
	return
}

func (sp *SessionPool) deleteSession(scope, key string) {
	k := sessionKey(scope, key)
	sp.sessionPool.Delete(k)
}

func (sp *SessionPool) NewUserSession(ctx context.Context, uid string) (sess *Session, ret context.Context, err error) {
	var loaded bool
	ret = ctx
	if sess, loaded = sp.createSession(scopeUser, uid); !loaded {
	}
	return
}

func (sp *SessionPool) GetUserSession(uid string) (sess *Session, ok bool) {
	return sp.loadSession(scopeUser, uid)
}

func (sp *SessionPool) DeleteUserSession(ctx context.Context, uid string) {
	sp.deleteSession(scopeUser, uid)
	return
}

func (sp *SessionPool) clearExpiredSession() {
	sp.sessionPool.Range(func(key, value interface{}) bool {
		v := value.(*Session)
		if !v.IsActive() {
			sp.sessionPool.Delete(key)
			slog.InfoContext(sp.ctx, "delete user session:%s due to expire ", key.(string))
		}
		return true
	})
}

// ForeachUID 遍历所有的在线用户
func (sp *SessionPool) ForeachUID(cb func(string, *Session) bool) {
	sp.sessionPool.Range(func(i, v interface{}) bool {
		uid := i.(string)
		sess := v.(*Session)
		uid = strings.TrimPrefix(uid, scopeUser+":")
		return cb(uid, sess)
	})
}

// SetOfflineHeart 多久没心跳，可以认为用户下线，默认60秒
func (sp *SessionPool) SetOfflineHeart(second int64) {
	if second > 0 {
		sp.offlineHeart = second * 1000
	}
}

func (sp *SessionPool) SetSessionExpire(millSecond int64) {
	if millSecond > 0 {
		sp.sessionExpireMillSecond = millSecond
	}
}
