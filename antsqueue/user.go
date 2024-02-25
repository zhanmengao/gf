package antsqueue

import (
	"time"
)

type UserRequestList struct {
	request    chan *TJob
	key        string
	qSize      int
	lastActive int64
}

func NewUserRequestList(key string, sz int) *UserRequestList {
	return &UserRequestList{
		key:   key,
		qSize: sz,
	}
}

func (u *UserRequestList) Init() {
	u.request = make(chan *TJob, u.qSize)
	u.lastActive = time.Now().Unix()
}

func (u *UserRequestList) Clear() {
	close(u.request)
}
