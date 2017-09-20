package job

import (
	"fmt"

	"github.com/kapmahc/fly/plugins/nut/health"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

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

func check() (interface{}, error) {
	err := open(func(_ *amqp.Channel) error {
		return nil
	})
	return "ok", err
}

func init() {
	health.Register("Redis", check)
}
