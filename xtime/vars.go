package xtime

import (
	"time"
)

var (
	offset         int64
	offsetDuration time.Duration
	loc            = time.Local
	// reg            = regexp.MustCompile(`[-: /]+`)
)

const (
	layoutTimeWithoutSeparator = "20060102150405"
	layoutDate                 = "20060102"
	layoutTime                 = "2006-01-02 15:04:05"
	layoutTimeDate             = "2006-01-02"
	msPerSecond                = 1000
	secondsPerDay              = 86400
	NaturalDay                 = msPerSecond * secondsPerDay
	NaturalWeekDay             = NaturalDay * 7
)
