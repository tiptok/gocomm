package models

import "log"

// 消息存储-发布
type MessageRepository interface {
	SaveMessage(message *Message) error
	FindNoPublishedStoredMessages() ([]*Message, error)
	FinishMessagesStatus(messageIds []int64, finishStatus int) error
}

// 消息存储-接收
type MessageReceiverRepository interface {
	ReceiveMessage(params map[string]interface{}) error
	ConfirmReceive(params map[string]interface{}) error
}

// 消费者
type Consumer interface {
	StartConsume() error
	WithTopicHandler(topic string, handler func(message interface{}) error)
	WithMessageReceiver(receiver MessageReceiverRepository)
}

// 生产者
type MessageProducer interface {
	Publish(messages []*Message, option map[string]interface{}) (*MessagePublishResult, error)
}

type LogInfo func(params ...interface{})

var DefaultLog LogInfo = func(params ...interface{}) {
	log.Println(params...)
}
