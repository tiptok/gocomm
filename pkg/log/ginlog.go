package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tiptok/gocomm/config"
	"io"
	"os"
	"strings"
)

func InitGinLog(conf config.Logger) {
	DefaultLog = newginlog(conf)

	formatMap = make(map[int]string, 5)
	formatMap[1] = v1
	formatMap[2] = v2
	formatMap[3] = v3
	formatMap[4] = v4
	formatMap[5] = v5
}

type ginLog struct {
	io.Writer
}

func newginlog(conf config.Logger) Log {
	f, _ := os.Create(conf.Filename)
	gin.DefaultWriter = io.MultiWriter(f)
	l := &ginLog{}
	return l
}

func (this *ginLog) Debug(args ...interface{}) {
	//this.log.Debug(args...)
	//beego.Debug(args...)
	fmt.Fprintf(gin.DefaultWriter, generateFmtStr(len(args)), args...)
}

func (this *ginLog) Info(args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, generateFmtStr(len(args)), args...)
}

func (this *ginLog) Warn(args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, generateFmtStr(len(args)), args...)
}

func (this *ginLog) Error(args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, generateFmtStr(len(args)), args...)
}

func (this *ginLog) Panic(args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, generateFmtStr(len(args)), args...)
}

func (this *ginLog) Fatal(args ...interface{}) {
	fmt.Fprintf(gin.DefaultWriter, generateFmtStr(len(args)), args...)
}

const (
	v1 = "%v \n"
	v2 = "%v %v \n"
	v3 = "%v %v %v \n"
	v4 = "%v %v %v %v \n"
	v5 = "%v %v %v %v %v \n"
)

var formatMap map[int]string

func generateFmtStr(n int) string {
	if v, ok := formatMap[n]; ok {
		return v
	}
	s := strings.Repeat("%v ", n) + "\n"
	return s
}
