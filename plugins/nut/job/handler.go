package job

import (
	log "github.com/sirupsen/logrus"
)

const (
	// PREFIX prefix
	PREFIX = "health://"
)

// Handler handler
type Handler func([]byte) (interface{}, error)

var handlers = make(map[string]Handler)

// Register register handler
func Register(n string, h Handler) {
	if _, ok := handlers[n]; ok {
		log.Warnf("task %s already exist, will override it", n)
	}
	handlers[n] = h
}
