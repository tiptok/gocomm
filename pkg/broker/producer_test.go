package broker

import (
	"github.com/tiptok/gocomm/identity/idgen"
	"github.com/tiptok/gocomm/pkg/broker/kafkax"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"testing"
	"time"
)

const KAFKA_HOSTS = "127.0.0.1:9092"

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

	producer := NewMessageProducer(nil, map[string]interface{}{"kafkaHosts": KAFKA_HOSTS})
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

//type PgMessageRepository struct {
//	transactionContext *transaction.TransactionContext
//}
//
//func (repository *PgMessageRepository) SaveMessage(message *models.Message) error {
//	sql := `insert into sys_message_produce (id,topic,value,msg_time,status)values(?,?,?,?,?)`
//	_, err := repository.transactionContext.PgDd.Exec(sql, message.Id, message.Topic, utils.JsonAssertString(message), message.MsgTime, int64(models.UnFinished))
//	return err
//}
//func (repository *PgMessageRepository) FindNoPublishedStoredMessages() ([]*models.Message, error) {
//	sql := `select value from sys_message_produce where status=?`
//	var values []string
//	_, e := repository.transactionContext.PgDd.Query(&values, sql, int64(models.UnFinished))
//	var messages = make([]*models.Message, 0)
//	if e != nil {
//		return messages, nil
//	}
//	for _, v := range values {
//		item := &models.Message{}
//		utils.JsonUnmarshal(v, item)
//		if item.Id != 0 {
//			messages = append(messages, item)
//		}
//	}
//	return messages, nil
//}
//func (repository *PgMessageRepository) FinishMessagesStatus(messageIds []int64, finishStatus int) error {
//	_, err := repository.transactionContext.PgDd.Exec("update sys_message_produce set status=? where id in (?)", finishStatus, pg.In(messageIds))
//	return err
//}
//func NewPgMessageRepository(transactionContext *transaction.TransactionContext) *PgMessageRepository {
//	return &PgMessageRepository{
//		transactionContext: transactionContext,
//	}
//}
