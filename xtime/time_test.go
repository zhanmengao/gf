package xtime

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSetTimeZone(t *testing.T) {
	err := SetTimeZone("America/New_York")
	if err != nil {
		panic(err)
	}
	now := Millisecond()
	str := TimeMSToString(now)
	log.Println("now:", str, "   ", now)
}

func TestSetTime(t *testing.T) {
	_, err := SetTime("2019-11-03 12:00:00")
	if err != nil {
		panic(err)
	}
	now := Millisecond()
	str := TimeMSToString(now)
	if str != "2019-11-03 12:00:00" {
		t.Fail()
	}
}

func TestDaysZeroTimeMs(t *testing.T) {
	ms := DaysZeroTimeMs(1)
	str := TimeMSToString(ms)
	if str != "2019-11-04 00:00:00" {
		t.Fail()
	}
}

func TestWeekDayZeroTimeMs(t *testing.T) {
	ms := WeekDayZeroTimeMs(time.Monday)
	str := TimeMSToString(ms)
	if str != "2019-11-04 00:00:00" {
		t.Fatalf("the date is :%s\n", str)
	}
}

func TestMonthDayZeroTime(t *testing.T) {
	ms := MonthDayZeroTime(3, 5)
	str := TimeMSToString(ms)
	if str != "2020-02-05 00:00:00" {
		t.Fatalf("the date is :%s\n", str)
	}
}

func TestNextPeriodWeekDayZeroTimeMs(t *testing.T) {
	SetTime("2019-11-01 12:00:01")
	//
	ms := NextPeriodWeekDayZeroTimeMs(time.Monday)
	str := TimeMSToString(ms)
	if str != "2019-11-04 00:00:00" {
		t.Fatalf("the date is :%s\n", str)
	}

	ms = NextPeriodWeekDayZeroTimeMs(time.Sunday)
	str = TimeMSToString(ms)
	if str != "2019-11-03 00:00:00" {
		t.Fatalf("the date is :%s\n", str)
	}

	ms = NextPeriodWeekDayZeroTimeMs(time.Friday)
	str = TimeMSToString(ms)
	if str != "2019-11-08 00:00:00" {
		t.Fatalf("the date is :%s\n", str)
	}

	ms = NextPeriodWeekDayZeroTimeMs(time.Saturday)
	str = TimeMSToString(ms)
	if str != "2019-11-02 00:00:00" {
		t.Fatalf("the date is :%s\n", str)
	}
}

func TestMSToLocalTimeString(t *testing.T) {
	str := "2019-11-20 00:00:00"
	SetTimeZone("UTC")
	ms := ParseToTimeUnixMS(str)

	str = MSToLocalTimeString(ms)
	if str != "2019-11-20 08:00:00" {
		t.Fatalf("str is :%s\n", str)
	}
}

func TestTimeZoneGMT(t *testing.T) {
	var data = []struct {
		tz  string
		gmt int
	}{
		{"Asia/Shanghai", 8},
		{"Canada/Pacific", -7},
		{"UTC", 0},
	}

	for i, d := range data {
		SetTimeZone(d.tz)
		gmt := TimeZoneGMT()
		if gmt != d.gmt {
			t.Fatalf("data %d check failed, gmt:%d -> %+v", i+1, gmt, d)
		}
	}
}

func TestInterval(t *testing.T) {
	begin := Unix() - 5*60*60
	end := Unix()
	fmt.Printf("begin = %v end = %v \n", time.Unix(begin, 0), time.Unix(end, 0))
	day := DayZeroInterval(begin, end)
	fmt.Println(day)
}
