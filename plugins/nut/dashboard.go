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

	var forum = []Link{
		{Href: "forum.Plugin.IndexArticles", Label: "forum.articles.index.title"},
		{Href: "forum.Plugin.IndexComments", Label: "forum.comments.index.title"},
	}
	if p.IsAdmin() {
		forum = append(forum, Link{Href: "forum.Plugin.IndexTags", Label: "forum.tags.index.title"})
	}
	p.AddDashboardMenu("forum.dashboard.title", forum...)

	p.AddDashboardMenu("survey.dashboard.title", Link{Href: "survey.Plugin.IndexForms", Label: "survey.forms.index.title"})

	if p.IsAdmin() {
		p.AddDashboardMenu(
			"nut.settings.title",
			Link{Href: "nut.Plugin.GetAdminSiteStatus", Label: "nut.admin.site.status.title"},
			Link{Href: "nut.Plugin.GetAdminSiteInfo", Label: "nut.admin.site.info.title"},
			Link{Href: "nut.Plugin.GetAdminSiteAuthor", Label: "nut.admin.site.author.title"},
			Link{Href: "nut.Plugin.GetAdminSiteSeo", Label: "nut.admin.site.seo.title"},
			Link{Href: "nut.Plugin.GetAdminSiteSMTP", Label: "nut.admin.site.smtp.title"},
			Link{Href: "nut.Plugin.IndexAdminLocales", Label: "nut.admin.locales.index.title"},
			Link{Href: "nut.Plugin.IndexAdminCards", Label: "nut.admin.cards.index.title"},
			Link{Href: "nut.Plugin.IndexAdminLinks", Label: "nut.admin.links.index.title"},
			Link{Href: "nut.Plugin.IndexAdminUsers", Label: "nut.admin.users.index.title"},
			Link{Href: "nut.Plugin.IndexAdminFriendLinks", Label: "nut.admin.friend-links.index.title"},
			Link{Href: "nut.Plugin.IndexLeaveWords", Label: "nut.leave-words.index.title"},
		)
	}

	p.Data["dashboard"] = p.dashboardMenus
	p.setFavicon()
	p.Layout = "layouts/dashboard/index.html"
}
