package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/kapmahc/fly/plugins/nut/app"
)

// Flush clear cache
func Flush() error {
	c := app.Redis().Get()
	defer c.Close()
	keys, err := redis.Values(c.Do("KEYS", PREFIX+"*"))
	if err == nil && len(keys) > 0 {
		_, err = c.Do("DEL", keys...)
	}
	return err
}
