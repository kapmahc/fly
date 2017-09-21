package nut

import (
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/kapmahc/fly/plugins/nut/i18n"
)

// Router get http router
func Router() *gin.Engine {
	return router
}

var router *gin.Engine

func openRouter() error {
	router = gin.Default()
	router.Use(i18n.DetectLocale)
	return nil
}

func init() {
	app.RegisterResource(openRouter)
}
