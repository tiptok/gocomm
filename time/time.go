package time

import (
	"fmt"
	"strconv"
	"time"
	"github.com/tiptok/gocomm/pkg/log"
)

//获取当前时间字符串,格式:"20170420133114" (2017-04-20 13:3114)
func GetTimeByYyyymmddhhmmss() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("20060102150405")
}

//获取当前时间字符串,格式:"0420133114" (2017-04-20 13:3114)
func GetTimeByhhmmss() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("150405")
}

func GetTimeByYyyymmddhhmm() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04")
}

// 获取当前日期前一天日期
func GetDateBeforeDay() string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, -1)
	logDay := yesTime.Format("20060102")
	return logDay
}

// 根据指定时间戳获取加减相应时间后的时间戳
func GetUnixTimeByUnix(timeUnix int64, years int, months int, days int) int64 {
	if timeUnix < 1 {
		return 0
	}
	tm := time.Unix(timeUnix, 0)
	return tm.AddDate(years, months, days).Unix()
}

//获取当前时间字符串,格式:"20170420" (2017-04-20)
func GetTimeByYyyymmdd() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("20060102")
}

func GetTimeByYyyymmdd2() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02")
}

//获取当前时间字符串,格式:"20170420" (2017-04-20)
func GetTimeByYyyymmddInt64() (int64, error) {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	nowDay, err := strconv.ParseInt(tm.Format("20060102"), 10, 64)
	if err != nil {
		return 0, err
	}
	return nowDay, nil
}

// 根据时间戳获取对应日期整数
func GetTDayByUnixTime(nowUnix int64) int64 {
	if nowUnix < 1 {
		return 0
	}
	tm := time.Unix(nowUnix, 0)
	nowDay, err := strconv.ParseInt(tm.Format("20060102"), 10, 64)
	if err != nil {
		log.Error(err)
		return 0
	}
	return nowDay
}

// 根据时间戳获取对应日期格式
func GetDiyTimeByUnixTime(nowUnix int64) string {
	if nowUnix < 1 {
		return ""
	}
	tm := time.Unix(nowUnix, 0)
	return tm.Format("2006/01/02")
}

// 根据时间戳获取对应月份整数
func GetMonthByUnixTime(nowUnix int64) int64 {
	if nowUnix < 1 {
		return 0
	}
	tm := time.Unix(nowUnix, 0)
	nowDay, err := strconv.ParseInt(tm.Format("200601"), 10, 64)
	if err != nil {
		log.Error(err)
		return 0
	}
	return nowDay
}

//获取当前日期(20170802)零点对应的Unix时间戳
func GetUnixTimeByYyyymmdd() int64 {
	timeStr := time.Now().Format("2006-01-02")

	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, err := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	if err != nil {
		log.Error(err)
		return 0
	}
	return t.Unix()
}

//获取指定时间戳下n天0点时间戳
func GetUnixTimeByNDayUnix(timeUnix int64, n int) int64 {
	timeUnix = GetUnixTimeByUnix(timeUnix, 0, 0, n)
	timeStr := time.Unix(timeUnix, 0).Format("2006-01-02")

	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, err := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	if err != nil {
		log.Error(err)
		return 0
	}
	return t.Unix()
}

//获取指定时间戳下n月0点时间戳
func GetUnixTimeByNMonthUnix(timeUnix int64, n int) int64 {
	timeUnix = GetUnixTimeByUnix(timeUnix, 0, n, 0)
	timeStr := time.Unix(timeUnix, 0).Format("2006-01-02")

	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t, err := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	if err != nil {
		log.Error(err)
		return 0
	}
	return t.Unix()
}

//获取指定时间下月份0点时间戳
func GetUnixTimeByMonthUnix(t time.Time)int64{
	year, month, _ := t.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return thisMonth.Unix()
}

// 获取制定时间戳是1970年1月1日开始的第几天
func GetDaythByTime(timeUnix int64) int64 {
	return (timeUnix+28800)/86400 + 1
}

// 获取上个月月初和月末的时间戳
func GetLastMonthStartAndEnd() (int64, int64) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, -1, 0).Unix()
	end := thisMonth.Unix() - 1
	return start, end
}

// 根据毫秒时间戳转换成20:18:23:3(20点28分23秒3毫秒)对应的整数(201823003)
func GetTimeNanoByNano(timeNano int64) int64 {
	tm := time.Unix(timeNano/1000, 0)
	str := fmt.Sprintf("%s%03d", tm.Format("150405"), timeNano%1000)
	n, _ := strconv.ParseInt(str, 10, 64)
	return n
}
