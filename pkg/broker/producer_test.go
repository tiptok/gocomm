package broker

import (
	"github.com/tiptok/gocomm/identity/idgen"
	"github.com/tiptok/gocomm/pkg/broker/kafkax"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"testing"
	"time"
)

const KAFKA_HOSTS = "106.52.15.41:9092"

//func TestNewMessageProducer(t *testing.T) {
//	var (
//		transactionContext = transaction.NewPGTransactionContext(pgDB.DB)
//		err                error
//	)
//
//	producer := NewMessageProducer(NewPgMessageRepository(transactionContext), map[string]interface{}{"kafkaHosts": KAFKA_HOSTS})
//	err = producer.PublishMessages([]*models.Message{
//		&models.Message{Id: idgen.Next(), Topic: "chat", MsgTime: time.Now().Unix(), Value: "hello world! tip tip!", FinishStatus: 0},
//	})
//
//	if err != nil {
//		return
//	}
//	time.Sleep(time.Second * 2)
//}

func TestNewMessageProducerNoRepository(t *testing.T) {
	var (
		err error
	)

	producer := NewMessageProducer(models.WithKafkaHost(KAFKA_HOSTS))
	err = producer.PublishMessages([]*models.Message{
		&models.Message{Id: idgen.Next(), Topic: "chat", MsgTime: time.Now().Unix(), Value: "hello world! tip tip!", FinishStatus: 0},
	})

	if err != nil {
		return
	}
	time.Sleep(time.Second * 2)
}

func TestProducer(t *testing.T) {
	producer, _ := kafkax.NewKafkaMessageProducer(KAFKA_HOSTS)
	for i := 0; i < 3; i++ {
		_, err := producer.Publish([]*models.Message{{Id: 1, Topic: "mmm_xcx_orders", MsgTime: time.Now().Unix(), Value: "hello ccc!"}}, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
}
