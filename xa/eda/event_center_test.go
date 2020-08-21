package eda

import (
	"fmt"
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
