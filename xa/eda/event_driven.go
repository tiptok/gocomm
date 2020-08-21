package eda

type (
	// 发布者
	EventPublisher interface {
		Subscribe(sub SubscriberHandler) error
		Publish(e Event) error
	}
	// 订阅处理者
	SubscriberHandler interface {
		HandleEvent(e Event) error
	}
	// 事件:定义触发的事件
	Event interface {
		// 事件类型
		EventType() string
	}
)

type HandleEvent func(e Event) error

// 基础发布者
type CommonEventPublisher struct {
	subscribers []SubscriberHandler
}

// Subscribe 订阅处理者
func (pub *CommonEventPublisher) Subscribe(sub ...SubscriberHandler) error {
	pub.subscribers = append(pub.subscribers, sub...)
	return nil
}

// Publish 发布事件
func (pub *CommonEventPublisher) Publish(e Event) error {
	for _, subscriber := range pub.subscribers {
		subscriber.HandleEvent(e)
	}
	return nil
}
