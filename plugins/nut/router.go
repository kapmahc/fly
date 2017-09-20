package nut

import (
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/fly/plugins/nut/app"
)

// Router get http router
func Router() *gin.Engine {
	return router
}

var router *gin.Engine

func openRouter() error {
	router = gin.Default()
	return nil
}

func init() {
	app.RegisterResource(openRouter)
}
