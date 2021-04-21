package broker

import (
	"github.com/Shopify/sarama"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/broker/local"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"testing"
)

func TestNewConsumer(t *testing.T) {
	consumer := NewConsumer(KAFKA_HOSTS, "0")
	consumer.WithMessageReceiver(local.NewPgMessageReceiverRepository(nil, nil))
	consumer.WithTopicHandler("mmm_xcx_orders", func(message interface{}) error {
		m, ok := message.(*sarama.Message)
		if !ok {
			return nil
		}
		if len(m.Value) > 0 {
			var msg models.Message
			common.JsonUnmarshal(string(m.Value), &msg)
			t.Log("handler message :", string(m.Value), msg.Id, msg.Topic, msg.Value)
		}
		return nil
	})
	consumer.StartConsume()
}
func TestNewConsumerNoRepository(t *testing.T) {
	consumer := NewConsumer(KAFKA_HOSTS, "mmm_orders")
	consumer.WithTopicHandler("mmm_xcx_orders", func(message interface{}) error {
		m, ok := message.(*sarama.Message)
		if !ok {
			return nil
		}
		if len(m.Value) > 0 {
			var msg models.Message
			common.JsonUnmarshal(string(m.Value), &msg)
			t.Log("handler message :", string(m.Value), msg.Id, msg.Topic, msg.Value)
		}
		return nil
	})
	consumer.StartConsume()
}
