package app

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

// Redis redis pool
func Redis() *redis.Pool {
	return _redis
}

var _redis *redis.Pool

func openRedis() error {
	_redis = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial(
				"tcp",
				fmt.Sprintf(
					"%s:%d",
					viper.GetString("redis.host"),
					viper.GetInt("redis.port"),
				),
			)
			if e != nil {
				return nil, e
			}
			if _, e = c.Do("SELECT", viper.GetInt("redis.db")); e != nil {
				c.Close()
				return nil, e
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func init() {
	RegisterResource(openRedis)
}
