package producer

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
)

var (
	DefaultExchangeType = "direct"
	DefaultReliable     = true
)

type JRabbitConfig struct {
	AMQP         string
	Exchange     string
	ExchangeType *string
	RoutingKey   *string
	Reliable     *bool
	Args         *amqp.Table
}

type JRabbitMQProducer struct {
	Client                             *amqp.Connection
	Exchange, ExchangeType, RoutingKey string
	Reliable                           bool
	AMQP                               string
	Args                               amqp.Table
	Channel                            *amqp.Channel
}

func NewJRabbitMQ(config JRabbitConfig) JRabbitMQProducer {
	client := JRabbitMQProducer{}
	connection, err := amqp.Dial(config.AMQP)
	if err != nil {
		log.Fatal("Dial: %s", err)
	}
	client.Client = connection
	client.Exchange = config.Exchange
	client.AMQP = config.AMQP
	client.ExchangeType = DefaultExchangeType
	if config.ExchangeType != nil {
		client.ExchangeType = *config.ExchangeType
	}

	client.Reliable = DefaultReliable
	if config.Reliable != nil {
		client.Reliable = *config.Reliable
	}
	//client.Channel, err = connection.Channel()
	//if err != nil {
	//    log.Fatal("Create channel fails : %s", err)
	//}
	if config.Args != nil {
		client.Args = *config.Args
	}
	return client
}

func init() {
	flag.Parse()

}
func initConnect(url string) (*amqp.Connection, error) {
	log.Infof("dialing %q", url)
	c := make(chan *amqp.Error)
	go func() {
		err := <-c
		fmt.Println("trying reconnect: " + err.Error())
		initConnect(url)
	}()
	connection, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Dial: %s", err)
	}

	connection.NotifyClose(c)
	log.Infof("got Connection, getting Channel")
	return connection, err
}

func (producer JRabbitMQProducer) PublishMessage(routingKey, body string, args amqp.Table) error {
	if producer.Client == nil || producer.Client.IsClosed() {
		connection, err := amqp.Dial(producer.AMQP)
		if err != nil {
			log.Fatal("Dial: %s", err)
		}
		producer.Client = connection
	}
	channel, err := producer.Client.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	defer channel.Close()

	//log.Infof("got Channel, declaring %q Exchange (%q)", producer.ExchangeType, producer.Exchange)
	if err := channel.ExchangeDeclare(
		producer.Exchange,     // name
		producer.ExchangeType, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // noWait
		producer.Args,         // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	// Reliable publisher confirms require confirm.select support from the
	// RabbitProducerClient.
	if producer.Reliable {
		//log.Infof("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}
		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer confirmOne(confirms)
	}
	//log.Infof("declared Exchange, publishing %dB body (%q)", len(body), body)
	if err = channel.Publish(
		producer.Exchange, // publish to an exchange
		routingKey,        // routing to 0 or more queues
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			Headers:         args,
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}
	return nil
}
