package health

import (
	"github.com/garyburd/redigo/redis"
	"github.com/kapmahc/fly/plugins/nut/app"
)

func postgresqlCheck() (interface{}, error) {
	// stats := app.DB().QueryOne(model, query, params)
	return app.DB().PoolStats(), nil
}

func check() (interface{}, error) {
	c := app.Redis().Get()
	defer c.Close()
	str, err := redis.String(c.Do("PING"))
	return str, err
}

func init() {
	Register("PostgreSQL", postgresqlCheck)
	Register("RabbitMQ", check)
}
