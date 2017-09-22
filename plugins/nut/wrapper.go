package nut

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/csrf"
)

const (
	// TITLE title
	TITLE = "title"

	// FLASH flash
	FLASH = "flash"
	// NOTICE flash notice
	NOTICE = "notice"
	// ERROR flash error
	ERROR = "error"
	// WARNING flash warning
	WARNING = "warning"
)

// FORM form handler
func FORM(to string, fm interface{}, fn func(*gin.Context, interface{}) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// https://github.com/gin-gonic/gin/issues/796
		// err := c.Bind(fm)
		err := binding.Form.Bind(c.Request, fm)
		if err == nil {
			err = fn(c, fm)
		}
		if err != nil {
			ss := sessions.Default(c)
			for _, v := range strings.Split(err.Error(), "\n") {
				ss.AddFlash(v, ERROR)
			}
			ss.Save()
		}

		c.Redirect(http.StatusFound, to)
	}
}

// HTML html handler
func HTML(t string, f func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			// TODO show error page
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ss := sessions.Default(c)
		fls := gin.H{}
		for _, k := range []string{NOTICE, WARNING, ERROR} {
			var msg []string
			for _, v := range ss.Flashes(k) {
				msg = append(msg, fmt.Sprintf("%v", v))
			}
			fls[k] = template.HTML(strings.Join(msg, "<br/>"))
		}
		ss.Save()

		c.Set(FLASH, fls)
		c.Set(csrf.TemplateTag, csrf.TemplateField(c.Request))

		c.HTML(http.StatusOK, t, c.Keys)
	}
}
