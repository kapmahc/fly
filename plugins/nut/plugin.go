package nut

import "github.com/astaxie/beego"

// Plugin controller
type Plugin struct {
	Controller
}

// GetHome home
// @router / [get]
func (p *Plugin) GetHome() {
	beego.Debug("##########")
	p.Data["Website"] = "beego.me"
	p.Data["Email"] = "astaxie@gmail.com"
	p.TplName = "index.tpl"
}
