package consumers

import (
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"strings"
	"sync/atomic"
	"time"

	"github.com/manucorporat/try"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	consumerTag  string // Name that consumer identifies itself to the server with
	uri          string // uri of the rabbitmq server
	exchange     string // exchange that we will bind to
	exchangeType string // topic, direct, etc...

	lastRecoverTime int64
	//track service current status
	currentStatus atomic.Value
}

func newConsumer(consumerTag, uri, exchange, exchangeType string) *Consumer {
	consumer := &Consumer{
		consumerTag:     fmt.Sprintf("%s", consumerTag),
		uri:             uri,
		exchange:        exchange,
		exchangeType:    exchangeType,
		done:            make(chan error),
		lastRecoverTime: time.Now().Unix(),
	}
	consumer.currentStatus.Store(true)
	return consumer
}

func RunConsumer(rabbitUri, consumerTag, exchange, exchangeType, queueName, routingKey string, mgo *database.MogoDB, config model.RankPointConfig) {
	consumer := newConsumer(consumerTag, rabbitUri, exchange, exchangeType)

	if err := consumer.Connect(); err != nil {
		fmt.Sprintf("[%s]connect error", consumerTag)
	}

	deliveries, err := consumer.AnnounceQueue(queueName, routingKey)
	if err != nil {
		log.Fatalf("[%s]Error when calling AnnounceQueue()", consumerTag)
	}
	consumer.Handle(deliveries, 2, queueName, routingKey, mgo, config)
}

func (c *Consumer) ReConnect(queueName, routingKey string, retryTime int) (<-chan amqp.Delivery, error) {
	c.Close()
	time.Sleep(time.Duration(config.Config.RabbitMq.TimeOutRetry) * time.Second)
	log.Debug("Try ReConnect with times: ", retryTime)

	if err := c.Connect(); err != nil {
		return nil, err
	}

	deliveries, err := c.AnnounceQueue(queueName, routingKey)
	if err != nil {
		return deliveries, errors.New("Couldn't connect")
	}
	return deliveries, nil
}

// Connect to RabbitMQ server
func (c *Consumer) Connect() error {

	var err error
	log.Debug("dialing: ", c.uri)
	c.conn, err = amqp.Dial(c.uri)

	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		// Waits here for the channel to be closed
		log.Info("closing: ", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Let Handle know it's not time to reconnect
		c.done <- errors.New("Channel Closed")
	}()

	log.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Debug("got Channel, declaring Exchange ", c.exchange)
	if err = c.channel.ExchangeDeclare(
		c.exchange,     // name of the exchange
		c.exchangeType, // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	return nil
}

// AnnounceQueue sets the queue that will be listened to for this
// connection...
func (c *Consumer) AnnounceQueue(queueName, routingKey string) (<-chan amqp.Delivery, error) {
	log.Info("declared Exchange, declaring Queue:", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}
	log.Info(fmt.Sprintf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, routingKey))
	err = c.channel.Qos(100, 0, false)
	if err != nil {
		return nil, fmt.Errorf("Error setting qos: %s", err)
	}
	arg := make(map[string]interface{})
	arg["x-dead-letter-exchange"] = config.Config.RabbitMq.Exchange
	for _, key := range strings.Split(routingKey, ",") {
		arg["x-dead-letter-routing-key"] = fmt.Sprintf("%s.dead", key)
		if err = c.channel.QueueBind(
			queue.Name, // name of the queue
			key,        // routingKey
			c.exchange, // sourceExchange
			false,      // noWait
			arg,        // arguments
		); err != nil {
			return nil, fmt.Errorf("Queue Bind: %s", err)
		}
	}

	log.Info("Queue bound to Exchange, starting Consume consumer tag:", c.consumerTag)
	deliveries, err := c.channel.Consume(
		queue.Name,    // name
		c.consumerTag, // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}
	return deliveries, nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
		c.channel = nil
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *Consumer) Handle(
	deliveries <-chan amqp.Delivery,
	threads int,
	queue string,
	routingKey string,
	mgo *database.MogoDB,
	config model.RankPointConfig) {

	var err error
	for {
		log.Info("Enter for busy loop with thread:", threads)
		for i := 0; i < threads; i++ {
			go func() {
				log.Info("Enter go with thread with deliveries", deliveries)
				for msg := range deliveries {
					ret := false
					try.This(func() {
						//log.Info(string(msg.Body))
						go func() {
							e := processing(msg.Body, mgo, config)
							if e != nil {
								log.Error(e)
							}
						}()

						ret = true
					}).Finally(func() {
						if ret == true {
							msg.Ack(false)
						} else {
							msg.Reject(false)
							// this really a litter dangerous. if the worker is panic very quickly,
							// it will ddos our sentry server......plz, add [retry-ttl] in header.
							//msg.Nack(false, true)
							c.currentStatus.Store(false)
						}
					}).Catch(func(e try.E) {
						log.Error(e)
					})
				}
			}()
		}

		if <-c.done != nil {
			c.currentStatus.Store(false)
			retryTime := 1
			for {
				deliveries, err = c.ReConnect(queue, routingKey, retryTime)
				if err != nil {
					log.Error("Reconnecting Error")
					retryTime += 1
				} else {
					break
				}
			}
		}
		log.Debug("Reconnected!!!")
	}
}
