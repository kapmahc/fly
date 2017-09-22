package nut

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/fly/plugins/nut/i18n"
)

func getInstall(c *gin.Context) error {
	lang := c.MustGet(i18n.LOCALE).(string)
	c.Set(TITLE, i18n.T(lang, "nut.install.title"))
	return nil
}

type fmInstall struct {
	Title                string `form:"title" binding:"required"`
	Subhead              string `form:"subhead" binding:"required"`
	Name                 string `form:"name" binding:"required"`
	Email                string `form:"name" binding:"required,email"`
	Password             string `form:"password" binding:"min=6"`
	PasswordConfirmation string `form:"passwordConfirmation" binding:"eqfield=Password"`
}

func postInstall(c *gin.Context, fm interface{}) error {
	log.Println("fuck")
	return nil
}
