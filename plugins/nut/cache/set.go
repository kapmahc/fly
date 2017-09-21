package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/kapmahc/fly/plugins/nut/app"
)

const (
	// PREFIX prefix
	PREFIX = "cache://"
)

//Set set cache item
func Set(key string, val interface{}, ttl time.Duration) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return err
	}

	c := app.Redis().Get()
	defer c.Close()
	_, err := c.Do("SET", PREFIX+key, buf.Bytes(), "EX", int(ttl/time.Second))
	return err
}
