package nut

// AddDashboardMenu add dashboard menu
func (p *Controller) AddDashboardMenu(label string, links ...Link) {
	var items []map[string]string
	for _, v := range links {
		items = append(items, map[string]string{
			"label": v.Label,
			"href":  v.Href,
		})
	}
	p.dashboardMenus = append(
		p.dashboardMenus,
		H{
			"label": label,
			"items": items,
		},
	)
}

// LayoutDashboard use dashboard layout
func (p *Controller) LayoutDashboard() {
	p.MustSignIn()
	p.dashboardMenus = make([]H, 0)
	p.AddDashboardMenu(
		"nut.self.title",
		Link{Href: "nut.Plugin.GetUsersLogs", Label: "nut.users.logs.title"},
		Link{Href: "nut.Plugin.GetUsersProfile", Label: "nut.users.profile.title"},
		Link{Href: "nut.Plugin.GetUsersChangePassword", Label: "nut.users.change-password.title"},
	)
	if p.IsAdmin() {
		p.AddDashboardMenu(
			"nut.dashboard.title",
			Link{Href: "nut.Plugin.GetUsersLogs", Label: "nut.users.logs.title"},
		)
	}

	p.Data["dashboard"] = p.dashboardMenus
	p.Layout = "layouts/dashboard/index.html"
}
