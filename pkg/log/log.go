package log

import "github.com/tiptok/gocomm/config"

type Log interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
}

var (
	DefaultLog Log = newConsoleLog()
)

func InitLog(conf config.Logger) {
	DefaultLog = newbeelog(conf)
}
func Debug(args ...interface{}) {
	DefaultLog.Debug(args...)
}
func Info(args ...interface{}) {
	DefaultLog.Info(args...)
}
func Warn(args ...interface{}) {
	DefaultLog.Warn(args...)
}
func Error(args ...interface{}) {
	DefaultLog.Error(args...)
}
func Panic(args ...interface{}) {
	DefaultLog.Panic(args...)
}
func Fatal(args ...interface{}) {
	DefaultLog.Fatal(args...)
}

func Logger() Log {
	return DefaultLog
}
