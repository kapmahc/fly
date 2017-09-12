package nut

// GetHome home
// @router / [get]
func (p *Plugin) GetHome() {
	p.LayoutApplication()
	p.TplName = "nut/home.html"
}
