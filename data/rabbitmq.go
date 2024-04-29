package data

import (
	"enclave_in_web3/utils"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"sync"
)

type RabbitMQ struct {
	Conn       *amqp.Connection
	Channel    *amqp.Channel
	Exchange   string
	QueueName  string
	RoutingKey string
	Url        string
}

func (mq *RabbitMQ) Close() {
	mq.Conn.Close()
	mq.Channel.Close()
}

var rabbitMqMgr *RabbitMqMgr

var ErrRabbitMqConfig = errors.New("rabbitmq config error")

var ErrRabbitMqUninitialized = errors.New("rabbitmq uninitialized")

func InitRabbitMqMgr() {
	rabbitMqMgr = newRabbitMqMgr(viper.Sub("data.rabbitmq"))
}

func ReleaseRabbitMqMgr() {
	if rabbitMqMgr != nil {
		rabbitMqMgr.Close()
		rabbitMqMgr = nil
	}
}

func GetRabbitMq(name string) (*RabbitMQ, error) {
	if rabbitMqMgr == nil {
		panic(ErrRabbitMqUninitialized)
	}

	return rabbitMqMgr.getRabbitMq(name)
}

func MustGetRabbitMq(name string) *RabbitMQ {
	if rabbitMqMgr == nil {
		panic(ErrRabbitMqUninitialized)
	}

	return rabbitMqMgr.mustGetRabbitMq(name)
}

func newRabbitMqMgr(conf *viper.Viper) *RabbitMqMgr {
	dbMgr := &RabbitMqMgr{
		mqMap:    make(map[string]*RabbitMQ),
		mutex:    &sync.Mutex{},
		mqConfig: conf,
	}
	return dbMgr
}

type RabbitMqMgr struct {
	mqMap    map[string]*RabbitMQ
	mutex    *sync.Mutex
	mqConfig *viper.Viper
}

func (mgr *RabbitMqMgr) getRabbitMq(name string) (*RabbitMQ, error) {
	config := mgr.mqConfig.Sub(name)
	if config == nil {
		return nil, ErrRabbitMqConfig
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mq, ok := mgr.mqMap[name]
	if ok {
		return mq, nil
	}

	mq, err := initRabbitMq(config, name)
	if err != nil {
		return nil, err
	}
	mgr.mqMap[name] = mq
	return mq, nil
}

func (mgr *RabbitMqMgr) mustGetRabbitMq(name string) *RabbitMQ {
	config := mgr.mqConfig.Sub(name)
	if config == nil {
		panic(ErrRabbitMqConfig)
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mq, ok := mgr.mqMap[name]
	if ok {
		return mq
	}

	mq, err := initRabbitMq(config, name)
	utils.CheckError(err)

	mgr.mqMap[name] = mq
	return mq
}

func (mgr *RabbitMqMgr) Close() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	for _, mq := range mgr.mqMap {
		mq.Close()
	}
	mgr.mqMap = make(map[string]*RabbitMQ)
}

type RabbitMqLogger struct {
}

func (logger *RabbitMqLogger) Print(values ...interface{}) {
	//glog.Info(gorm.LogFormatter(values...)...)
	fmt.Println(gorm.LogFormatter(values...)...)
}

func initRabbitMq(config *viper.Viper, name string) (*RabbitMQ, error) {
	url := config.GetString("url")
	exchange := config.GetString("exchange")
	queue := config.GetString("queue")
	routingKey := config.GetString("routing-key")
	rabbitMq := &RabbitMQ{
		Url:        url,
		Exchange:   exchange,
		QueueName:  queue,
		RoutingKey: routingKey,
	}

	var err error
	rabbitMq.Conn, err = amqp.Dial(rabbitMq.Url)
	if err != nil {
		return rabbitMq, err
	}

	rabbitMq.Channel, err = rabbitMq.Conn.Channel()
	if err != nil {
		return rabbitMq, err
	}

	_, err = rabbitMq.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return rabbitMq, err
	}

	err = rabbitMq.Channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return rabbitMq, err
	}

	err = rabbitMq.Channel.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		return rabbitMq, err
	}

	glog.Infof("%s rabbitmq: url: %s, exchange: %s, queue: %s, routing key: %s",
		name, url, exchange, queue, routingKey)

	return rabbitMq, nil
}

func (mq *RabbitMQ) Publish(contentType string, content string) error {
	return mq.Channel.Publish(
		mq.Exchange,
		mq.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType, //"text/plain",
			Body:        []byte(content),
		},
	)
}

func (mq *RabbitMQ) PublishPlain(content string) error {
	return mq.Channel.Publish(
		mq.Exchange,
		mq.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		},
	)
}

func (mq *RabbitMQ) PublishJson(content string) error {
	return mq.Channel.Publish(
		mq.Exchange,
		mq.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(content),
		},
	)
}

func (mq *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	return mq.Channel.Consume(
		mq.QueueName,
		"gallium",
		true,
		false,
		false,
		true,
		nil,
	)
}
