package kafkax

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/identity/idgen"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/sync/task"
	"strings"
	"sync"
	"time"
)

type SaramaConsumer struct {
	ready             chan bool
	messageHandlerMap map[string]func(message interface{}) error
	//Logger            log.Logger
	kafkaHosts string
	groupId    string
	topicMiss  map[string]string //记录未被消费的topic
	receiver   models.MessageReceiverRepository
	//version    string
	option        *models.MessageOptions
	ConsumerRetry *ConsumerRetry
}

func (consumer *SaramaConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}
func (consumer *SaramaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
func (consumer *SaramaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var err error
	for message := range claim.Messages() {
		err = consumer.messageProcess(message)
		session.MarkMessage(message, "")
		if err != nil {
			// message retry
			if consumer.option.ConsumeRetryOption.Enable {
				consumer.ConsumerRetry.StoreRetryMessage(message)
			}
			continue
		}
	}
	return err
}

func (consumer *SaramaConsumer) messageProcess(message *sarama.ConsumerMessage) error {
	var err error
	log.Debug(fmt.Sprintf("【kafka】 receive message  topic = %s offset = %v value = %v", message.Topic, message.Offset, string(message.Value)))
	handler, ok := consumer.messageHandlerMap[message.Topic]
	var msg = &models.Message{}
	common.JsonUnmarshal(string(message.Value), msg)

	if e := consumer.messageReceiveBefore(message, msg.Id); e != nil {
		ok = false
		log.Error(e)
	}
	if !ok {
		return nil
	}
	var handlerMsg interface{} = msg
	if consumer.option.HandlerOriginalMessageFlag {
		handlerMsg = message
	}
	if err = handler(handlerMsg); err != nil {
		log.Error("【kafka】 message process error topic =", message.Topic, "offset:", message.Offset, err)
		return err
	}
	consumer.messageReceiveAfter(msg.Id)
	return nil
}
func (consumer *SaramaConsumer) messageReceiveBefore(message *sarama.ConsumerMessage, msgId int64) error {
	if consumer.receiver == nil {
		return nil
	}

	var params = make(map[string]interface{})
	var err error
	_, ok := consumer.messageHandlerMap[message.Topic]
	if !ok {
		params["status"] = models.Ignore
		_, topicMiss := consumer.topicMiss[message.Topic]
		if !topicMiss {
			fmt.Printf("topic:[%v] has not consumer handler", message.Topic)
		}
		return nil
	}

	_, err = consumer.storeMessage(params, message, msgId)
	if err != nil {
		ok = false
		//log.Println("ConsumeClaim:", err)
	}
	return err
}
func (consumer *SaramaConsumer) messageReceiveAfter(msgId int64) {
	if consumer.receiver == nil {
		return
	}
	consumer.finishMessage(map[string]interface{}{"id": msgId})
}

func (consumer *SaramaConsumer) storeMessage(params map[string]interface{}, message *sarama.ConsumerMessage, msgId int64) (id int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
		}
	}()
	id = idgen.Next()
	params = make(map[string]interface{})
	params["id"] = message.Offset
	params["topic"] = message.Topic
	params["partition"] = message.Partition
	params["offset"] = message.Offset
	params["key"] = string(message.Key)
	params["value"] = string(message.Value)
	params["msg_time"] = message.Timestamp.Unix()
	params["create_at"] = time.Now().Unix()
	params["status"] = models.UnFinished //0:未完成 1:已完成 2：未命中
	if msgId > 0 {
		params["id"] = msgId
	}
	err = consumer.receiver.ReceiveMessage(params)
	return
}
func (consumer *SaramaConsumer) finishMessage(params map[string]interface{}) error {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
		}
	}()
	consumer.receiver.ConfirmReceive(params)
	return nil
}

func (consumer *SaramaConsumer) StartConsume() error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	//config.Consumer.Offsets.AutoCommit.Enable = false
	config.Version = sarama.V0_10_2_1
	if len(consumer.option.Version) > 0 {
		if v, e := sarama.ParseKafkaVersion(consumer.option.Version); e != nil {
			return e
		} else {
			config.Version = v
		}
	}
	brokerList := strings.Split(consumer.kafkaHosts, ",")
	consumerGroup, err := sarama.NewConsumerGroup(brokerList, consumer.groupId, config)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			var topics []string
			for key := range consumer.messageHandlerMap {
				topics = append(topics, key)
			}
			if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
				log.Error(err.Error())
				return
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	<-consumer.ready
	log.Info("Sarama consumer up and running")

	consumer.StartExtraWork()
	select {
	case <-ctx.Done():
		log.Info("Sarama consumer : context cancelled")
	}
	cancel()
	wg.Wait()
	if err := consumerGroup.Close(); err != nil {
		return err
	}
	return nil
}
func (consumer *SaramaConsumer) StartExtraWork() {
	if consumer.option.ConsumeRetryOption.Enable {
		consumer.ConsumerRetry = NewConsumerRetry(consumer.option.ConsumeRetryOption, consumer)
		consumer.ConsumerRetry.task.Start()
		log.Info(fmt.Sprintf("ConsumerRetry start  maxtime:%v duration:%vs", consumer.option.ConsumeRetryOption.MaxRetryTime, consumer.option.ConsumeRetryOption.NextRetryTimeSpan))
	}
}
func (consumer *SaramaConsumer) WithTopicHandler(topic string, handler func(message interface{}) error) { //*sarama.ConsumerMessage
	consumer.messageHandlerMap[topic] = handler
}
func (consumer *SaramaConsumer) WithMessageReceiver(receiver models.MessageReceiverRepository) {
	consumer.receiver = receiver
}

func NewSaramaConsumer(kafkaHosts string, groupId string, options ...models.MessageOption) models.Consumer {
	var option = models.NewMessageOptions(options...)
	return &SaramaConsumer{
		kafkaHosts:        kafkaHosts,
		groupId:           groupId,
		topicMiss:         make(map[string]string),
		messageHandlerMap: make(map[string]func(message interface{}) error),
		ready:             make(chan bool),
		option:            option,
	}
}

type ConsumerRetry struct {
	option *models.ConsumeRetryOption
	// 容器
	store models.MessageStore
	//retryHandler interface{}
	exit chan int
	// periodic task
	task *task.Periodic
	// consumer
	consumer *SaramaConsumer
}

func (retry *ConsumerRetry) ConsumeRetryMessage(consumer *SaramaConsumer) {
	retry.task.Start()
}

func (retry *ConsumerRetry) StoreRetryMessage(message *sarama.ConsumerMessage) error {
	retryMessage := &models.RetryMessage{
		Message:       message,
		RetryTime:     1,
		MaxRetryTime:  retry.option.MaxRetryTime,
		NextRetryTime: time.Now().Add(time.Duration(retry.option.NextRetryTimeSpan)).Unix(),
	}
	return retry.store.StoreMessage(retryMessage)
}

func (retry *ConsumerRetry) processMessage() error {
	defer func() {
		if p := recover(); p != nil {
			log.Warn(p)
		}
	}()
	messages, err := retry.store.GetMessage()
	if err != nil {
		log.Error(err)
		return err
	}
	if len(messages) == 0 {
		return nil
	}
	for _, m := range messages {
		if m.RetryTime > retry.option.MaxRetryTime {
			continue
		}
		err := retry.consumer.messageProcess(m.Message)
		if err != nil && m.RetryTime < retry.option.MaxRetryTime {
			m.RetryTime++
			m.NextRetryTime = time.Now().Add(time.Second * time.Duration(retry.option.NextRetryTimeSpan)).Unix()
			retry.store.StoreMessage(m)
		}
	}
	return nil
}

func NewConsumerRetry(option *models.ConsumeRetryOption, consumer *SaramaConsumer) *ConsumerRetry {
	retry := &ConsumerRetry{
		option:   option,
		store:    option.Store,
		consumer: consumer,
	}
	retry.task = task.NewPeriodic(time.Second*time.Duration(option.NextRetryTimeSpan), retry.processMessage)
	return retry
}
