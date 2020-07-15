package log

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tiptok/gocomm/config"
	"path/filepath"
	"strconv"
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
	ilv, err := strconv.Atoi(conf.Level)
	if err != nil {
		ilv = logs.LevelDebug
	}
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
