package health

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/fly/plugins/nut/app"
	log "github.com/sirupsen/logrus"
)

// Check run check
func Check() {
	for n, h := range handlers {
		log.Info("do helath check ", n)
		now := time.Now()
		v, e := h()
		d := time.Now().Sub(now)

		if err := saveResult(n, d, v, e); err != nil {
			log.Error(err)
		}
		log.Info(d, v, e)
	}
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
