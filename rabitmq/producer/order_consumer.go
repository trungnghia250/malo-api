package producer

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/trungnghia250/malo-api/config"
	"github.com/trungnghia250/malo-api/database"
	"github.com/trungnghia250/malo-api/service/model"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	rabbitCf = config.Config.RabbitMq
	producer JRabbitMQProducer
)

func init() {
	rabbitConfig := JRabbitConfig{
		AMQP:     rabbitCf.AMQP,
		Exchange: rabbitCf.Exchange,
	}
	producer = NewJRabbitMQ(rabbitConfig)

}
func OrderStream(col *mongo.Collection) {
	var pipeline []echo.Map
	ctx := context.Background()
	// stream func to handle event
	streamer := func(cs *mongo.ChangeStream) {
		for cs.Next(ctx) {
			changeDoc := new(model.Mgostream)
			if err := cs.Decode(changeDoc); err != nil {
				log.Error("decode error:", err)
				break
			}
			announceEvent(changeDoc)
		}
	}
	database.ChangeStream(col, pipeline, streamer)
}

func announceEvent(docChange *model.Mgostream) {
	go func() {
		b, _ := json.Marshal(docChange)
		producer.PublishMessage(docChange.OPType, string(b), nil)
	}()

}
