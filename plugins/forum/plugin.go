package forum

import "github.com/kapmahc/fly/plugins/nut"

// Plugin controller
type Plugin struct {
	nut.Controller
}

// GetHome home
// @router / [get]
func (p *Plugin) GetHome() {
	p.TplName = "index.tpl"
}
