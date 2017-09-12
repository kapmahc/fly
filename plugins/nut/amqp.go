package nut

import (
	"fmt"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

var (
	_jobber    *Jobber
	jobberOnce sync.Once
)

// JOBBER get a jobber
func JOBBER() *Jobber {
	jobberOnce.Do(func() {
		_jobber = &Jobber{
			queue:    beego.BConfig.ServerName,
			source:   beego.AppConfig.String("amqpsource"),
			handlers: make(map[string]JobHandler),
		}
	})
	return _jobber
}

// JobHandler job handler
type JobHandler func(body []byte) error

// Jobber jobber
type Jobber struct {
	queue    string
	source   string
	handlers map[string]JobHandler
}

// Register registe job handler
func (p *Jobber) Register(n string, h JobHandler) {
	if _, ok := p.handlers[n]; ok {
		beego.Warn("already have", n, ", will override it")
	}
	p.handlers[n] = h
}

// Receive receive jobs
func (p *Jobber) Receive(consumer string) error {
	beego.Info("waiting for messages, to exit press CTRL+C")
	return p.open(func(ch *amqp.Channel) error {
		if err := ch.Qos(1, 0, false); err != nil {
			return err
		}
		qu, err := ch.QueueDeclare(p.queue, true, false, false, false, nil)
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
			beego.Info("receive message", d.MessageId, "@", d.Type)
			hnd, ok := p.handlers[d.Type]
			if !ok {
				return fmt.Errorf("unknown message type %s", d.Type)
			}
			if err := hnd(d.Body); err != nil {
				return err
			}
			beego.Info("done", d.MessageId, time.Now().Sub(now))
		}
		return nil
	})
}

// Send send job
func (p *Jobber) Send(pri uint8, typ string, body []byte) error {
	return p.open(func(ch *amqp.Channel) error {
		qu, err := ch.QueueDeclare(p.queue, true, false, false, false, nil)
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

func (p *Jobber) open(f func(*amqp.Channel) error) error {
	conn, err := amqp.Dial(p.source)
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
