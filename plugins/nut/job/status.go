package job

import (
	"github.com/garyburd/redigo/redis"
	"github.com/kapmahc/fly/plugins/nut/app"
)

// Status status
func Status(start, stop int) map[string][]string {
	c := app.Redis().Get()
	defer c.Close()
	items := make(map[string][]string)
	for n := range handlers {
		if val, err := redis.Strings(c.Do("LRANGE", PREFIX+n, start, stop)); err == nil {
			items[n] = val
		}
	}
	return items
}
