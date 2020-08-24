// eda 基于事件驱动实现架构
package eda

import (
	"fmt"
	"reflect"
	"sync"
)

// 事件中心
// 通过统一的事件中心,发布订阅事件,
// 实现基于事件驱动的处理方式,解耦代码
type EventCenter interface {
	// 发布
	Publish(event Event) error
	// 注册订阅者
	RegisterSubscribe(event Event, sub SubscriberHandler) error
	// 撤销订阅者
	DeregisterSubscribe(event Event, sub SubscriberHandler) error
}

var defaultRegisterCenter = newDefaultRegisterCenter()

func newDefaultRegisterCenter() EventCenter {
	return &DefaultRegisterCenter{
		subscribers: make(map[string]map[string]HandleEvent),
	}
}

type DefaultRegisterCenter struct {
	subscribers map[string]map[string]HandleEvent
	regMutex    sync.Mutex
}

func (r *DefaultRegisterCenter) Publish(event Event) error {
	handleList, ok := r.subscribers[event.EventType()]
	if !ok || len(handleList) == 0 {
		return nil
	}
	for i := range handleList {
		err := handleList[i](event)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *DefaultRegisterCenter) RegisterSubscribe(event Event, sub SubscriberHandler) error {
	r.regMutex.Lock()
	defer r.regMutex.Unlock()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	key := reflect.ValueOf(sub).Elem().String()
	if handleList, ok := r.subscribers[event.EventType()]; ok {

		if _, ok := handleList[key]; ok {
			return fmt.Errorf("%v has register", key)
		}
		handleList[key] = sub.HandleEvent
		return nil
	}
	var newHandleList = make(map[string]HandleEvent)
	newHandleList[key] = sub.HandleEvent
	r.subscribers[event.EventType()] = newHandleList
	return nil
}
func (r *DefaultRegisterCenter) DeregisterSubscribe(event Event, sub SubscriberHandler) error {
	r.regMutex.Lock()
	defer r.regMutex.Unlock()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	//defer fmt.Println("unlock:",event.EventType())
	//fmt.Println("lock:",event.EventType())
	key := reflect.ValueOf(sub).Elem().String()
	if handleList, ok := r.subscribers[event.EventType()]; ok {
		if _, ok := handleList[key]; ok {
			delete(handleList, key)
		}
		if len(handleList) == 0 {
			delete(r.subscribers, event.EventType())
		}
		return nil
	}
	return nil
}

// 事件中心发布者
type EventCenterPublisher struct {
	e EventCenter
}

func NewEventCenterPublisher(v EventCenter) EventCenterPublisher {
	return EventCenterPublisher{e: v}
}

// Subscribe 订阅处理者
func (pub *EventCenterPublisher) Subscribe(sub ...SubscriberHandler) error {
	return nil
}

// Publish 发布事件
func (pub *EventCenterPublisher) Publish(e Event) error {
	if pub.e == nil {
		return defaultRegisterCenter.Publish(e)
	}
	return pub.e.Publish(e)
}

/*************************对外接口***************************/
func Publish(event Event) error {
	return defaultRegisterCenter.Publish(event)
}
func RegisterSubscribe(event Event, sub SubscriberHandler) error {
	if defaultRegisterCenter == nil {
		return fmt.Errorf("default register center is nil")
	}
	return defaultRegisterCenter.RegisterSubscribe(event, sub)
}
func DeregisterSubscribe(event Event, sub SubscriberHandler) error {
	if defaultRegisterCenter == nil {
		return fmt.Errorf("default register center is nil")
	}
	return defaultRegisterCenter.DeregisterSubscribe(event, sub)
}
