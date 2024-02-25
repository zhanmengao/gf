package xtime

import (
	"time"
)

func Reset() {
	offset = 0
	offsetDuration = 0
}

func SetTime(tf string) (t time.Time, err error) {
	t, err = ParseTime(tf)
	if err != nil {
		return
	}
	offset = (t.Unix() - time.Now().Unix()) * msPerSecond
	offsetDuration = time.Duration(offset) * time.Millisecond
	return
}

func Offset() int64 {
	return offset
}
