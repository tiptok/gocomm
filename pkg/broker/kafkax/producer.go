package kafkax

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/identity/idgen"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"log"
	"strings"
	"time"
)

// sarame kafka 消息生产
type KafkaMessageProducer struct {
	KafkaHosts string
	LogInfo    models.LogInfo
	producer   sarama.SyncProducer
}

func NewKafkaMessageProducer(host string) (*KafkaMessageProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Retry.Max = 10
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Version = sarama.V0_11_0_0
	brokerList := strings.Split(host, ",")
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}
	return &KafkaMessageProducer{KafkaHosts: host, LogInfo: models.DefaultLog, producer: producer}, nil
}

// 同步发送
func (engine *KafkaMessageProducer) Publish(messages []*models.Message, option map[string]interface{}) (*models.MessagePublishResult, error) {
	if engine.producer == nil {
		return nil, fmt.Errorf("producer haven`t set up")
	}
	var successMessageIds []int64
	var errMessageIds []int64
	for _, message := range messages {
		if value, err := json.Marshal(message); err == nil {
			msg := &sarama.ProducerMessage{
				Topic:     message.Topic,
				Value:     sarama.StringEncoder(value),
				Timestamp: time.Now(),
			}
			partition, offset, err := engine.producer.SendMessage(msg)
			if err != nil {
				errMessageIds = append(errMessageIds, message.Id)
				log.Println(err)
			} else {
				successMessageIds = append(successMessageIds, message.Id)
				var append = make(map[string]interface{})
				append["topic"] = message.Topic
				append["partition"] = partition
				append["offset"] = offset
				log.Println("kafka消息发送", append)
			}
		}
	}
	return &models.MessagePublishResult{SuccessMessageIds: successMessageIds, ErrorMessageIds: errMessageIds}, nil
}

// 消息调度器
type MessageDispatcher struct {
	notifications     chan struct{}
	messageChan       chan *models.Message
	dispatchTicker    *time.Ticker
	messageRepository models.MessagePublisherRepository
	producer          models.MessageProducer
}

func (dispatcher *MessageDispatcher) MessagePublishedNotice() error {
	time.Sleep(time.Second * 2)
	dispatcher.notifications <- struct{}{}
	return nil
}

func (dispatcher *MessageDispatcher) MessagePublish(messages []*models.Message) error {
	for i := range messages {
		dispatcher.messageChan <- messages[i]
	}
	return nil
}

// go dispatcher.Dispatch() 启动一个独立协程
func (dispatcher *MessageDispatcher) Dispatch() {
	for {
		select {
		case <-dispatcher.dispatchTicker.C:
			go func(dispatcher *MessageDispatcher) {
				dispatcher.notifications <- struct{}{}
			}(dispatcher)
		case <-dispatcher.notifications:
			if dispatcher.messageRepository == nil {
				continue
			}
			messages, _ := dispatcher.messageRepository.FindNoPublishedStoredMessages()
			var messagesInProcessIds []int64
			for i := range messages {
				messagesInProcessIds = append(messagesInProcessIds, messages[i].Id)
			}
			if messages != nil && len(messages) > 0 {
				dispatcher.messageRepository.FinishMessagesStatus(messagesInProcessIds, int(models.InProcess))

				reuslt, err := dispatcher.producer.Publish(messages, nil)
				if err == nil && len(reuslt.SuccessMessageIds) > 0 {
					dispatcher.messageRepository.FinishMessagesStatus(reuslt.SuccessMessageIds, int(models.Finished))
				}
				//发送失败的消息ID列表 更新状态 进行中->未开始
				if len(reuslt.ErrorMessageIds) > 0 {
					dispatcher.messageRepository.FinishMessagesStatus(reuslt.ErrorMessageIds, int(models.UnFinished))
				}
			}
		case msg := <-dispatcher.messageChan:
			dispatcher.producer.Publish([]*models.Message{msg}, nil)
		}
	}
}

type MessageDirector struct {
	messageRepository models.MessagePublisherRepository
	dispatcher        *MessageDispatcher
}

func (d *MessageDirector) Publish(topic string, originalMessages []interface{}, options ...models.MessageOption) error {
	var message []*models.Message
	for i := range originalMessages {
		m := originalMessages[i]
		message = append(message, &models.Message{
			Id:           idgen.Next(),
			Topic:        topic,
			Value:        common.JsonAssertString(m),
			MsgTime:      time.Now().Unix(),
			FinishStatus: 0,
		})
	}
	return d.PublishMessages(message, options...)
}

func (d *MessageDirector) PublishMessages(messages []*models.Message, options ...models.MessageOption) error {
	var option = models.NewMessageOptions(options...)
	if d.dispatcher == nil {
		return fmt.Errorf("dispatcher还没有启动")
	}
	if d.messageRepository == nil {
		d.dispatcher.MessagePublish(messages)
		return nil
	}
	for _, message := range messages {
		if option.MessageRepository == nil {
			break
		}
		if err := option.MessageRepository.SaveMessage(message); err != nil {
			return err
		}
	}
	go d.dispatcher.MessagePublishedNotice()
	return nil
}

// 消息发布器
// options["kafkaHosts"]="localhost:9092"
// options["timeInterval"]=time.Second*60*5
func NewMessageDirector(options ...models.MessageOption) *MessageDirector {
	var option = models.NewMessageOptions(options...)

	dispatcher := &MessageDispatcher{
		notifications:     make(chan struct{}),
		messageRepository: option.MessageRepository,
		messageChan:       make(chan *models.Message, 100),
	}
	var err error
	dispatcher.producer, err = NewKafkaMessageProducer(option.KafkaHost)
	if err != nil {
		log.Println(err)
		return nil
	}

	dispatcher.dispatchTicker = time.NewTicker(option.Interval)
	go dispatcher.Dispatch()

	return &MessageDirector{
		messageRepository: option.MessageRepository,
		dispatcher:        dispatcher,
	}
}
