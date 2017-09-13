package nut

// GetHome home
// @router / [get]
func (p *Plugin) GetHome() {
	p.LayoutApplication()
	p.Data[TITLE] = Tr(p.Locale(), "nut.home.title")
	p.TplName = "nut/home.html"
}

// GetDashboard dashboard panel
// @router /dashboard [get]
func (p *Plugin) GetDashboard() {
	p.MustSignIn()
	p.LayoutDashboard()
	p.Data[TITLE] = Tr(p.Locale(), "nut.dashboard.title")
	p.TplName = "nut/dashboard.html"
}
