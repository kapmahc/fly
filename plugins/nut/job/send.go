package job

import (
	"time"

	"github.com/google/uuid"
	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/streadway/amqp"
)

// Send send job
func Send(pri uint8, typ string, body []byte) error {
	return open(func(ch *amqp.Channel) error {
		qu, err := ch.QueueDeclare(app.Name(), true, false, false, false, nil)
		if err != nil {
			return err
		}

		return ch.Publish("", qu.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			MessageId:    uuid.New().String(),
			Priority:     pri,
			Body:         body,
			Timestamp:    time.Now(),
			Type:         typ,
		})
	})
}
