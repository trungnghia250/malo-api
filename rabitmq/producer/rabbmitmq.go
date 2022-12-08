package producer

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
	"github.com/trungnghia250/malo-api/config"
)

var (
	uri          = flag.String("uri", config.Config.RabbitMq.AMQP, "AMQP URI")
	exchangeName = flag.String("exchange", config.Config.RabbitMq.Exchange, "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", config.Config.RabbitMq.RoutingKey, "AMQP routing key")
	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.Parse()
}

func PublishMessage(message, key string) {
	if err := publish(*uri, *exchangeName, *exchangeType, key, message, *reliable); err != nil {
		log.Fatalf("%s", err)
	}

}
func publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {
	//log.Infof("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	//log.Infof("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	//log.Infof("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if reliable {
		//log.Infof("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}
		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer confirmOne(confirms)
	}
	//log.Infof("declared Exchange, publishing %dB body (%q)", len(body), body)
	if err = channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
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

func confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; !confirmed.Ack {
		log.Infof("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
