package log

import (
	"github.com/beego/beego/v2/core/logs"
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
	logs.Debug("", args...)
}

func (this *beegoLog) Info(args ...interface{}) {
	logs.Info("", args...)
}

func (this *beegoLog) Warn(args ...interface{}) {
	logs.Warn("", args...)
}

func (this *beegoLog) Error(args ...interface{}) {
	logs.Error("", args...)
}

func (this *beegoLog) Panic(args ...interface{}) {
	logs.Error("", args...)
}

func (this *beegoLog) Fatal(args ...interface{}) {
	logs.Error("", args...)
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
