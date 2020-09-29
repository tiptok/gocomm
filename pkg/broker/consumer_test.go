package broker

import (
	"github.com/Shopify/sarama"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"testing"
)

//func TestNewConsumer(t *testing.T) {
//	consumer := NewConsumer(constant.KAFKA_HOSTS, "0")
//	consumer.WithMessageReceiver(NewPgMessageReceiverRepository(transaction.NewPGTransactionContext(pg.DB)))
//	consumer.WithTopicHandler("mmm_xcx_orders", func(message interface{}) error {
//		m, ok := message.(*sarama.Message)
//		if !ok {
//			return nil
//		}
//		if len(m.Value) > 0 {
//			var msg models.Message
//			utils.JsonUnmarshal(string(m.Value), &msg)
//			t.Log("handler message :", string(m.Value), msg.Id, msg.Topic, msg.Value)
//		}
//		return nil
//	})
//	consumer.StartConsume()
//}
func TestNewConsumerNoRepository(t *testing.T) {
	consumer := NewConsumer(KAFKA_HOSTS, "0")
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

//type PgMessageReceiverRepository struct {
//	//transactionContext *transaction.TransactionContext
//}
//
////func NewPgMessageReceiverRepository(transactionContext *transaction.TransactionContext) *PgMessageReceiverRepository {
////	return &PgMessageReceiverRepository{
////		transactionContext: transactionContext,
////	}
////}
//func (repository *PgMessageReceiverRepository) ReceiveMessage(params map[string]interface{}) error {
//	//var num int
//	//checkSql := `select count(0) from sys_message_consume where "offset" =? and topic=?`
//	//_, err := repository.transactionContext.PgDd.Query(&num, checkSql, params["offset"], params["topic"])
//	//if err != nil {
//	//	return err
//	//}
//	//if num > 0 {
//	//	return fmt.Errorf("receive repeate message [%v]", params)
//	//}
//	//
//	//sql := `insert into sys_message_consume(topic,partition,"offset",key,value,msg_time,create_at,status)values(?,?,?,?,?,?,?,?)`
//	//_, err = repository.transactionContext.PgDd.Exec(sql, params["topic"], params["partition"], params["offset"], params["key"], params["value"], params["msg_time"], params["create_at"], params["status"])
//	//return err
//	return nil
//}
//func (repository *PgMessageReceiverRepository) ConfirmReceive(params map[string]interface{}) error {
//	//log.Println(params)
//	//_, err := repository.transactionContext.PgDd.Exec(`update sys_message_consume set status=? where "offset" =? and topic=?`, int(models.Finished), params["offset"], params["topic"])
//	//return err
//	return nil
//}
