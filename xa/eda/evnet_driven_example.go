package eda

import "fmt"

type ServiceCar struct {
	CommonEventPublisher
}

func (s *ServiceCar) Drive(speed int, wheel int) error {
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

// domain->event
// 定义领域事件
const DRIVE_CAR_EVENT = "drive:car:event"

type DriveEvent struct {
	// 速度
	Speed int
	// 转弯角度 -90-0-90
	Wheel int
}

func (e DriveEvent) EventType() string {
	return DRIVE_CAR_EVENT
}

// domain->service
// 定义领域服务
type Car interface {
	EventPublisher
	Drive() error
}

//application/event/subscriber
// 定义事件订阅者

// 汽车引擎
type CarEngine struct{}

func (s *CarEngine) HandleEvent(event Event) error {
	switch event.EventType() {
	case DRIVE_CAR_EVENT:
		d := event.(DriveEvent)
		fmt.Println(s.String(), "速度:", d.Speed, "转向:", d.Wheel)
		break
	default:
		break
	}
	return nil
}
func (s *CarEngine) String() string {
	return "引擎"
}

// 汽车油箱
type CarFuelTank struct{}

func (s *CarFuelTank) HandleEvent(event Event) error {
	switch event.EventType() {
	case DRIVE_CAR_EVENT:
		d := event.(DriveEvent)
		fmt.Println(s.String(), "输出油:", (float64(d.Speed)*1.0)/100.0, "KM/L")
		break
	default:
		break
	}
	return nil
}
func (s *CarFuelTank) String() string {
	return "油箱"
}
