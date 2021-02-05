package broker

import (
	"github.com/tiptok/gocomm/pkg/broker/kafkax"
	"github.com/tiptok/gocomm/pkg/broker/models"
)

//新消费者-消费组
func NewConsumer(kafkaHosts string, groupId string, options ...models.MessageOption) models.Consumer {
	return kafkax.NewSaramaConsumer(kafkaHosts, groupId, options...)
}
