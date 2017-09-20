package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/garyburd/redigo/redis"
	"github.com/kapmahc/fly/plugins/nut/app"
)

//Get get from cache
func Get(key string, val interface{}) error {
	c := app.Redis().Get()
	defer c.Close()
	bys, err := redis.Bytes(c.Do("GET", PREFIX+key))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(bys)
	return dec.Decode(val)
}
