package xtime

import (
	"time"
)

//Time format
const (
	YYYYMMDDHHMMSS = "20060102150405"
	YYYYMMDD       = "20060102"
	HHMMSS         = "150405"

	YYYYMMDDHHMMSSRFC3339 = "2006-01-02 15:04:05"
	YYYYMMDDRFC3339       = "2006-01-02"
)

type XTime time.Time

func (o XTime) Format(format string) string {
	return time.Time(o).Format(format)
}

//DayBefore 当前日期的前n天 零点时间
// n<0  当前时间之前
// n>0  当前时间之后
func (o XTime) DayBefore(x int) XTime {
	t := time.Time(o)
	yesTime := t.AddDate(0, 0, x)
	y, m, d := yesTime.Date()
	return XTime(time.Date(y, m, d, 0, 0, 0, 0, time.Local))
}

//DayBefore 当前日期的前n月 零点时间
// n<0  当前时间之前
// n>0  当前时间之后
func (o XTime) MonthBefore(x int) XTime {
	t := time.Time(o)
	yesTime := t.AddDate(0, x, 0)
	y, m, _ := yesTime.Date()
	return XTime(time.Date(y, m, 1, 0, 0, 0, 0, time.Local))
}

//昨天
func (o XTime) Yesterday() XTime {
	return o.DayBefore(-1)
}
