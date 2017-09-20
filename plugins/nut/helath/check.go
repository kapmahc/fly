package helath

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

// Check run check
func Check(f func(time.Duration, string, error)) {
	for n, h := range handlers {
		now := time.Now()
		log.Info("do helath check", n)
		s := ""
		v, e := h()
		if e == nil {
			var b []byte
			b, e = json.Marshal(v)
			if e == nil {
				s = string(b)
			}
		}
		f(time.Now().Sub(now), s, e)
	}
}
