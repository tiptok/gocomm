package task

import (
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"log"
	"strings"
	"testing"
	"time"
)

func Test_OnSuccess(t *testing.T){
	work :=func()error{
		log.Println("do work in")
		return errors.New("do work error")
	}
	afterwork:= func()error{
		log.Println("after work")
		return nil
	}
	f :=OnSuccess(work,afterwork)
	err := f()
	if err!=nil{
		log.Println(err)
	}
}

func Test_ExecuteParallel(t *testing.T){
	err :=Run(context.Background(),
		func() error {
			time.Sleep(time.Microsecond*300)
			return errors.New("T1")
		},
		func()error{
			time.Sleep(time.Microsecond*500)
			return errors.New("T2")
		})
	if r:=cmp.Diff(err.Error(),"T1");r!=""{
		t.Error(r)
	}
}

func Test_ExecuteParallelContextCancel(t *testing.T){
	ctx,cancel :=context.WithCancel(context.Background())
	err :=Run(ctx,
		func() error {
			time.Sleep(time.Microsecond*3000)
			return errors.New("T1")
		},
		func()error{
			time.Sleep(time.Microsecond*5000)
			return errors.New("T2")
		},
		func()error{
			time.Sleep(time.Microsecond*1000)
			cancel()
			return nil
		})
	errStr := err.Error()
	if strings.Contains(errStr, "canceled") {
		t.Error("expected error string to contain 'canceled', but actually not: ", errStr)
	}
}

func BenchmarkExecuteOne(b *testing.B){
	noop:=func()error{
		return nil
	}
	for i:=0;i<b.N;i++{
		Run(context.Background(),noop)
	}
}

func BenchmarkExecuteTwo(b *testing.B){
	noop:=func()error{
		return nil
	}
	for i:=0;i<b.N;i++{
		Run(context.Background(),noop,noop)
	}
}