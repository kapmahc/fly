package nut

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

// HTTPMethodFilter parse _method for put and delete
func HTTPMethodFilter(ctx *context.Context) {
	if ctx.Input.Query("_method") != "" && ctx.Input.IsPost() {
		beego.Debug(ctx.Request.Header)
		ctx.Request.Method = strings.ToUpper(ctx.Input.Query("_method"))
	}
}

func init() {
	// beego.InsertFilter("*", beego.BeforeRouter, HTTPMethodFilter)
}
