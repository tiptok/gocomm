package eda

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

type ECServiceCar struct {
	EventCenterPublisher
}

func (s *ECServiceCar) Drive(speed int, wheel int) error {
	e := DriveEvent{
		Speed: speed,
		Wheel: wheel,
	}
	if e := s.Publish(e); e != nil {
		fmt.Println(e)
		return e
	}
	return nil
}

func TestEventCenterExample(t *testing.T) {
	carSvr := &ECServiceCar{}

	//统一注册事件订阅
	RegisterSubscribe(DriveEvent{}, &CarEngine{})
	RegisterSubscribe(DriveEvent{}, &CarFuelTank{})

	//统一发布事件
	carSvr.Drive(50, 90)
	carSvr.Drive(100, 80)
	carSvr.Drive(120, 90)

	//统一解除事件订阅
	DeregisterSubscribe(DriveEvent{}, &CarFuelTank{})
	carSvr.Drive(50, 90)
	carSvr.Drive(50, 100)
}

func TestDefaultRegisterCenter_DeregisterSubscribe(t *testing.T) {
	center := &DefaultRegisterCenter{
		subscribers: make(map[string]map[string]HandleEvent),
	}
	for i := 0; i < 100; i++ {
		center.RegisterSubscribe(CustomerEvent(i), &CarEngine{})
	}
	if len(center.subscribers) != 100 {
		t.Fatal("error register number")
	}
	var wg = new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		value := i
		go func() {
			defer wg.Done()
			center.DeregisterSubscribe(CustomerEvent(value), &CarEngine{})
		}()
	}
	wg.Wait()
	if len(center.subscribers) > 0 {
		t.Fatal("error deregister number", len(center.subscribers))
	}
}

func TestDefaultRegisterCenter_RegisterSubscribe(t *testing.T) {
	var wg = new(sync.WaitGroup)
	center := &DefaultRegisterCenter{
		subscribers: make(map[string]map[string]HandleEvent),
	}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		value := i
		go func() {
			defer wg.Done()
			center.RegisterSubscribe(CustomerEvent(value), &CarEngine{})
		}()
	}
	wg.Wait()
	if len(center.subscribers) != 100 {
		t.Fatal("error register number", len(center.subscribers))
	}
}

func BenchmarkDefaultRegisterCenter(b *testing.B) {
	center := &DefaultRegisterCenter{
		subscribers: make(map[string]map[string]HandleEvent),
	}
	for i := 0; i < b.N; i++ {
		v := i
		e := center.RegisterSubscribe(CustomerEvent(v), &CarEngine{})
		if e != nil {
			b.Fatal(e)
		}
		e = center.DeregisterSubscribe(CustomerEvent(v), &CarEngine{})
		if e != nil {
			b.Fatal(e)
		}
	}
}

type CustomerEvent int

func (c CustomerEvent) EventType() string {
	return "customer" + strconv.Itoa(int(c))
}
