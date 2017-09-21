package app

import (
	"github.com/gin-gonic/gin"
)

// MountFunc http  mount func
type MountFunc func(*gin.Engine) error

var _mount []MountFunc

// RegisterMount register mount func
func RegisterMount(args ...MountFunc) {
	_mount = append(_mount, args...)
}

// Router return a http.Router
func Router() (*gin.Engine, error) {
	if IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	rt := gin.Default()
	for _, h := range _mount {
		if err := h(rt); err != nil {
			return nil, err
		}
	}
	return rt, nil
}
