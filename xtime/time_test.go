package xtime

import (
	"testing"
	"time"
)

func TestXTime_DayBefore(t *testing.T) {
	xtime, e := time.Parse(YYYYMMDD, "20200715")
	if e != nil {
		t.Fatal(e)
	}
	xdayBefore := XTime(xtime).DayBefore(-1)
	if xdayBefore.Format(YYYYMMDD) != "20200714" {
		t.Fatal("xtime error", xdayBefore.Format(YYYYMMDD))
	}
	xdayNext := XTime(xtime).DayBefore(1)
	if xdayNext.Format(YYYYMMDD) != "20200716" {
		t.Fatal("xtime error", xdayNext.Format(YYYYMMDD))
	}
	xmonthBefore := XTime(xtime).MonthBefore(-1)
	if xmonthBefore.Format(YYYYMMDD) != "20200601" {
		t.Fatal("xtime error", xmonthBefore.Format(YYYYMMDD))
	}
	xmonthNext := XTime(xtime).MonthBefore(1)
	if xmonthNext.Format(YYYYMMDD) != "20200801" {
		t.Fatal("xtime error", xmonthNext.Format(YYYYMMDD))
	}
}

func BenchmarkXTime_DayBefore(b *testing.B) {
	xtime, e := time.Parse(YYYYMMDD, "20200715")
	if e != nil {
		b.Fatal(e)
	}
	var xdayBefore XTime
	for i := 0; i < b.N; i++ {
		xdayBefore = XTime(xtime).DayBefore(0)
		if xdayBefore.Format(YYYYMMDD) != "20200715" {
			b.Fatal("xtime error", xdayBefore.Format(YYYYMMDD))
		}
	}
}
