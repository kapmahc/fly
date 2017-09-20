package job

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kapmahc/fly/plugins/nut/app"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

// Receive receive jobs
func Receive(consumer string) error {
	log.Info("waiting for messages, to exit press CTRL+C")
	return open(func(ch *amqp.Channel) error {
		if err := ch.Qos(1, 0, false); err != nil {
			return err
		}
		qu, err := ch.QueueDeclare(app.Name(), true, false, false, false, nil)
		if err != nil {
			return err
		}
		msgs, err := ch.Consume(qu.Name, consumer, false, false, false, false, nil)
		if err != nil {
			return err
		}
		for d := range msgs {
			d.Ack(false)
			now := time.Now()
			log.Info("receive message", d.MessageId, "@", d.Type)
			hnd, ok := handlers[d.Type]
			if !ok {
				return fmt.Errorf("unknown message type %s", d.Type)
			}
			if err := hnd(d.Body); err != nil {
				return err
			}
			log.Info("done", d.MessageId, time.Now().Sub(now))
		}
		return nil
	})
}

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

func open(f func(*amqp.Channel) error) error {
	arg := viper.GetStringMap("rabbitmq")
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		arg["user"].(string),
		arg["password"].(string),
		arg["host"].(string),
		arg["port"].(int64),
		arg["virtual"].(string),
	))
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return f(ch)
}
