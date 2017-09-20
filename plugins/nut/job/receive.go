package job

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/fly/plugins/nut/app"
	log "github.com/sirupsen/logrus"
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
			log.Info("receive message", d.MessageId, "@", d.Type)
			hnd, ok := handlers[d.Type]
			if !ok {
				return fmt.Errorf("unknown message type %s", d.Type)
			}

			now := time.Now()
			res, err := hnd(d.Body)
			saveResult(d.Type, time.Now().Sub(now), res, err)

			if err != nil {
				return err
			}
			log.Info("done", d.MessageId)
		}
		return nil
	})
}

func saveResult(n string, d time.Duration, v interface{}, e error) error {
	buf, err := json.Marshal(gin.H{
		"name":    n,
		"spend":   d,
		"result":  v,
		"error":   e,
		"created": time.Now(),
	})
	if err != nil {
		return err
	}
	c := app.Redis().Get()
	defer c.Close()
	_, err = c.Do("LPUSH", PREFIX+n, buf)
	return err
}
