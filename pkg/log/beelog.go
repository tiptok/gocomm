package log

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tiptok/gocomm/config"
	"path/filepath"
	"strings"
)

type beegoLog struct {
	log *logs.BeeLogger
}

func newbeelog(conf config.Logger) Log {
	filename := `{"filename":"` + filepath.ToSlash(conf.Filename) + `"}`

	l := &beegoLog{
		log: logs.GetBeeLogger(),
	}
	l.log.SetLogger(logs.AdapterFile, filename)
	ilv := beegoLogLevelAdapter(conf.Level)
	l.log.SetLevel(ilv)
	l.log.EnableFuncCallDepth(true)
	l.log.SetLogFuncCallDepth(6)
	return l
}

func (this *beegoLog) Debug(args ...interface{}) {
	//this.log.Debug(args...)
	beego.Debug(args...)
}

func (this *beegoLog) Info(args ...interface{}) {
	beego.Info(args...)
}

func (this *beegoLog) Warn(args ...interface{}) {
	beego.Warn(args...)
}

func (this *beegoLog) Error(args ...interface{}) {
	beego.Error(args...)
}

func (this *beegoLog) Panic(args ...interface{}) {
	beego.Error(args...)
}

func (this *beegoLog) Fatal(args ...interface{}) {
	beego.Error(args...)
}

func beegoLogLevelAdapter(logLevel string) int {
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		return logs.LevelDebug
	case "INFO":
		return logs.LevelInformational
	case "WARN":
		return logs.LevelWarning
	case "ERROR":
		return logs.LevelError
	}
	return logs.LevelDebug
}
