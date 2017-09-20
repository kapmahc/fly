package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/kapmahc/fly/plugins/nut/app"
)

// List return cache keys
func List() ([]string, error) {
	c := app.Redis().Get()
	defer c.Close()
	keys, err := redis.Strings(c.Do("KEYS", PREFIX+"*"))
	if err != nil {
		return nil, err
	}
	for i := range keys {
		keys[i] = keys[i][len(PREFIX):]
	}
	return keys, nil
}
