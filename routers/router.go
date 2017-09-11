package routers

import (
	"github.com/astaxie/beego"
	"github.com/kapmahc/fly/plugins/erp"
	"github.com/kapmahc/fly/plugins/forum"
	"github.com/kapmahc/fly/plugins/mall"
	"github.com/kapmahc/fly/plugins/nut"
	"github.com/kapmahc/fly/plugins/ops/mail"
	"github.com/kapmahc/fly/plugins/ops/vpn"
	"github.com/kapmahc/fly/plugins/pos"
	"github.com/kapmahc/fly/plugins/reading"
	"github.com/kapmahc/fly/plugins/survey"
)

func init() {
	beego.Include(&nut.Plugin{})
	for k, v := range map[string]beego.ControllerInterface{
		"/forum":    &forum.Plugin{},
		"/survey":   &survey.Plugin{},
		"/reading":  &reading.Plugin{},
		"/mall":     &mall.Plugin{},
		"/pos":      &pos.Plugin{},
		"/erp":      &erp.Plugin{},
		"/ops/mail": &mail.Plugin{},
		"/ops/vpn":  &vpn.Plugin{},
	} {
		beego.AddNamespace(beego.NewNamespace(k, beego.NSInclude(v)))
	}
}
