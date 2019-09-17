package task

import (
	"github.com/tiptok/gocomm/common"
	"log"
	"testing"
	"time"
)

func TestPeriodic(t *testing.T){
	count:=0
	task :=NewPeriodic(time.Second*2,func()error{
		count++
		log.Println("current count:",count)
		return nil
	})
	common.Must(task.Start())
	time.Sleep(time.Second * 5)
	common.Must(task.Close())
	log.Println("Count:",count)
	common.Must(task.Start())
	time.Sleep(time.Second*5)
	log.Println("Count:",count)
	common.Must(task.Close())
}
