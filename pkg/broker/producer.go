package broker

import (
	"github.com/tiptok/gocomm/pkg/broker/kafkax"
	"github.com/tiptok/gocomm/pkg/broker/models"
)

// 消息发布器
// options["kafkaHosts"]="localhost:9092"
// options["timeInterval"]=time.Second*60*5
func NewMessageProducer(messageRepository models.MessageRepository, options map[string]interface{}) *kafkax.MessageDirector {
	dispatcher := kafkax.NewMessageDirector(messageRepository, options)
	return dispatcher
}
