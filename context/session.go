package framework

import (
	"context"
	"github.com/zhanmengao/gf/util"
	"github.com/zhanmengao/gf/xtime"
	"strings"
	"sync"
)

const sessionKeyUID = "sess_uid"
const (
	SessionKeyClient    = "base_sess_client"
	SessionKeyRouteAddr = "base_session_route_%s_addr"
	SessionKeyAes       = "base_session_aes_key"
)

// Session 用来表示用户在服务器内生存时间之内的
type Session struct {
	createTime     int64
	lastActiveTime int64
	lastHeartTime  int64
	data           map[string]interface{}
	ttlMap         map[string]int64
	sync.RWMutex
	srvList []struct {
		srcName string
		srvAddr string
	}
	offlineHeart            *int64
	sessionExpireMillSecond *int64
}

func newSession(offlineHeart *int64, sessionExpireMillSecond *int64, uid string) *Session {
	now := xtime.Millisecond()
	s := &Session{
		createTime:              now,
		lastActiveTime:          now,
		lastHeartTime:           now,
		data:                    make(map[string]interface{}),
		offlineHeart:            offlineHeart,
		sessionExpireMillSecond: sessionExpireMillSecond,
		ttlMap:                  make(map[string]int64, 0),
	}
	s.setUID(uid)
	return s
}

func sessionKey(scope, key string) string {
	return scope + ":" + key
}

func sessionUID(sessKey string) string {
	ss := strings.Split(sessKey, ":")
	if len(ss) >= 2 {
		return ss[1]
	}
	return ""
}

func (sess *Session) GetCreateTime() int64 {
	sess.RLock()
	defer sess.RUnlock()
	return sess.createTime
}

func (sess *Session) GetSessValue(key string) (value interface{}) {
	sess.RLock()
	defer sess.RUnlock()
	return sess.getSessValue(key)
}

func (sess *Session) getSessValue(key string) (value interface{}) {
	var ok bool
	if value, ok = sess.data[key]; ok {
		return value
	}
	return nil
}

func (sess *Session) SetSessValue(key string, value interface{}) {
	sess.Lock()
	defer sess.Unlock()
	sess.setSessValue(key, value)
}

func (sess *Session) ModifySessValue(cb func(data map[string]interface{})) {
	sess.Lock()
	defer sess.Unlock()
	util.SafeFunc(context.Background(), func(ctx context.Context) {
		cb(sess.data)
	})
}

func (sess *Session) setSessValue(key string, value interface{}) {
	sess.data[key] = value
}

func (sess *Session) GetLastActiveTime() int64 {
	sess.RLock()
	defer sess.RUnlock()
	return sess.lastActiveTime
}

func (sess *Session) IsOnline() bool {
	sess.RLock()
	defer sess.RUnlock()
	nt := xtime.Millisecond()
	// 1分钟没有心跳就认为离线
	return nt-sess.lastHeartTime <= *sess.offlineHeart
}

func (sess *Session) UpdateHeartBeatTime() {
	sess.Lock()
	defer sess.Unlock()
	sess.lastHeartTime = xtime.Millisecond()
}

func (sess *Session) UpdateActiveTime() {
	sess.Lock()
	defer sess.Unlock()
	sess.lastActiveTime = xtime.Millisecond()
}

func (sess *Session) IsActive() bool {
	sess.RLock()
	defer sess.RUnlock()
	return xtime.Millisecond()-sess.lastActiveTime <= *sess.sessionExpireMillSecond
}

func (sess *Session) UpdateSrvList(srvName, srvAddr string) {
	sess.Lock()
	defer sess.Unlock()
	for i := range sess.srvList {
		v := &sess.srvList[i]
		if v.srcName == srvName && v.srvAddr != srvAddr {
			v.srvAddr = srvAddr
			return
		} else if v.srcName == srvName && v.srvAddr == srvAddr {
			return
		}
	}
	sess.srvList = append(sess.srvList, struct {
		srcName string
		srvAddr string
	}{srcName: srvName, srvAddr: srvAddr})
}

func (sess *Session) GetAddr(srvName string) string {
	sess.RLock()
	defer sess.RUnlock()
	list := sess.srvList
	for i := range list {
		if list[i].srcName == srvName {
			return list[i].srvAddr
		}
	}
	return ""
}

func (sess *Session) GetUID() string {
	if iUID := sess.GetSessValue(sessionKeyUID); iUID != nil {
		return iUID.(string)
	}
	return ""
}

func (sess *Session) setUID(uid string) {
	sess.SetSessValue(sessionKeyUID, uid)
}

func (sess *Session) GetFromCache(ctx context.Context, key string) (data interface{}, ok bool) {
	sess.RLock()
	defer sess.RUnlock()
	//ttl
	if expire, exist := sess.ttlMap[key]; exist {
		if expire > 0 && expire < xtime.Unix() {
			//过期了，不返回数据
			return
		}
	}
	//从sess读
	data = sess.getSessValue(key)
	if data != nil {
		ok = true
	}
	return
}
func (sess *Session) SetToCache(ctx context.Context, key string, val interface{}) {
	sess.Lock()
	defer sess.Unlock()
	//删除超时
	delete(sess.ttlMap, key)
	sess.setSessValue(key, val)
	return
}
func (sess *Session) SetToCacheTTL(ctx context.Context, key string, val interface{}, ttl int) {
	sess.Lock()
	defer sess.Unlock()
	//设置超时
	sess.ttlMap[key] = xtime.Unix() + int64(ttl)
	sess.setSessValue(key, val)
	return
}
