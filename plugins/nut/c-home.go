package nut

// GetHome home
// @router / [get]
func (p *Plugin) GetHome() {
	p.LayoutApplication()
	var tpl string
	if err := Get("home.tpl", &tpl); err != nil {
		tpl = "offcanvas"
	}
	p.Data[TITLE] = Tr(p.Locale(), "nut.home.title")
	p.TplName = "nut/home/" + tpl + ".html"
}

// GetDashboard dashboard panel
// @router /dashboard [get]
func (p *Plugin) GetDashboard() {
	p.LayoutDashboard()
	p.Data[TITLE] = Tr(p.Locale(), "nut.dashboard.title")
	p.TplName = "nut/dashboard.html"
}
