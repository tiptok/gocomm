package log

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type GzLog struct {
}

func InitGzLog(conf logx.LogConf) error {
	DefaultLog = &GzLog{}
	return logx.SetUp(conf)
}

func (this *GzLog) Debug(args ...interface{}) {
	fmt.Println(args...)
}

func (this *GzLog) Info(args ...interface{}) {
	logx.Info(args...)
}

func (this *GzLog) Warn(args ...interface{}) {
	logx.Error(args...)
}

func (this *GzLog) Error(args ...interface{}) {
	logx.Error(args...)
}

func (this *GzLog) Panic(args ...interface{}) {
	logx.Error(args...)
}

func (this *GzLog) Fatal(args ...interface{}) {
	logx.Error(args...)
}
