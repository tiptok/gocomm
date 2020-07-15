package log

import (
	"log"
)

type ConsoleLog struct {
}

func newConsoleLog() Log {
	return &ConsoleLog{}
}

func (this *ConsoleLog) Debug(args ...interface{}) {
	//this.log.Debug(args...)
	log.Println(args...)
}

func (this *ConsoleLog) Info(args ...interface{}) {
	log.Println(args...)
}

func (this *ConsoleLog) Warn(args ...interface{}) {
	log.Println(args...)
}

func (this *ConsoleLog) Error(args ...interface{}) {
	log.Println(args...)
}

func (this *ConsoleLog) Panic(args ...interface{}) {
	log.Println(args...)
}

func (this *ConsoleLog) Fatal(args ...interface{}) {
	log.Println(args...)
}
