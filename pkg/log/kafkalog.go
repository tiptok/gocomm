package log

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/Shopify/sarama"
	"github.com/beego/beego/v2/core/logs"
)

const loggerName = "kafkalog"
const MaxMessageSize = 500

var (
	ErrorInvalidKafkaConfig = fmt.Errorf("kafka config invalid")
	ErrorMessageSize        = fmt.Errorf("massage size over limit:%v", MaxMessageSize)
)

type KafkaLogger struct {
	done     chan struct{}
	config   *KafkaConfig
	msg      chan string
	size     int32
	closed   int32
	producer sarama.SyncProducer
}
type KafkaConfig struct {
	Topic   string   `json:"topic"`
	Level   int      `json:"level"`
	Key     string   `json:"key"`
	Addrs   []string `json:"addrs"`
	MaxSize int
}

func InitKafkaLogger(config KafkaConfig) (err error) {
	logs.Register(loggerName, NewKafkaLogger)
	jsondata, _ := json.Marshal(config)
	logs.SetLogger(loggerName, string(jsondata))
	return
}

/*
	实现 logger 接口
*/
func NewKafkaLogger() logs.Logger {
	log := &KafkaLogger{
		msg: make(chan string, MaxMessageSize),
	}

	go log.ConsumeMsg()
	return log
}
func (log *KafkaLogger) Init(configstr string) error {
	var (
		c   *KafkaConfig
		err error
	)
	if err = json.Unmarshal([]byte(configstr), &c); err != nil {
		return err
	}
	log.config = c
	if len(c.Topic) == 0 || len(c.Addrs) == 0 {
		return ErrorInvalidKafkaConfig
	}
	if log.config.MaxSize == 0 {
		log.config.MaxSize = MaxMessageSize
	}

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_11_0_0
	if log.producer, err = sarama.NewSyncProducer(c.Addrs, config); err != nil {
		return err
	}
	return nil
}
func (log *KafkaLogger) WriteMsg(lm *logs.LogMsg) error {
	//var when time.Time = lm.When
	var msg string = lm.Msg
	var level int = lm.Level
	if log.size >= MaxMessageSize {
		return ErrorMessageSize
	}
	if log.closed == 1 { //关闭停止接收
		return nil
	}
	if log.config.Level != 0 && level > log.config.Level {
		return nil
	}
	log.msg <- msg
	atomic.AddInt32(&log.size, 1)
	return nil
}
func (log *KafkaLogger) Destroy() {
	close(log.msg)
	log.producer.Close()
}
func (log *KafkaLogger) Flush() {
	close(log.done)
	atomic.CompareAndSwapInt32(&log.closed, 0, 1)
	//for msg,ok:=range log.msg{
	//  //send msg to kafka
	//}
}
func (log *KafkaLogger) SetFormatter(f logs.LogFormatter) {

}

func (log *KafkaLogger) ConsumeMsg() {
	for {
		select {
		case <-log.done:
			return
		case m, ok := <-log.msg:
			atomic.AddInt32(&log.size, -1)
			if ok {
				if _, _, err := log.producer.SendMessage(&sarama.ProducerMessage{
					Topic: log.config.Topic,
					Key:   sarama.ByteEncoder(log.config.Key),
					Value: sarama.ByteEncoder(m),
				}); err != nil {
					//TODO: err handler
				}
			}
		}
	}
}
