package xtime

import "time"

func SetTimeZone(tz string) (err error) {
	switch tz {
	case "Local", "":
		loc = time.Local
	default:
		lz, e := time.LoadLocation(tz)
		if e != nil {
			err = e
			return
		}
		loc = lz
	}
	return
}

func ResetTimeZone() {
	loc = time.Local
}

// TimeZoneGMT 返回时区的偏移
// eg: GMT+8 return 8, GMT-4 return -4
func TimeZoneGMT() int {
	t := time.Now().In(loc)
	_, offset := t.Zone()
	return offset / 3600
}
