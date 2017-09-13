package nut

import "github.com/astaxie/beego/orm"

// GetUsersProfile users'profile
// @router /users/profile [get]
func (p *Plugin) GetUsersProfile() {
	p.LayoutDashboard()
	p.Data[TITLE] = Tr(p.Locale(), "nut.users.profile.title")
	p.TplName = "nut/users/profile.html"
}

// GetUsersChangePassword change user password
// @router /users/change-password [get]
func (p *Plugin) GetUsersChangePassword() {
	p.LayoutDashboard()
	p.Data[TITLE] = Tr(p.Locale(), "nut.users.change-password.title")
	p.TplName = "nut/users/change-password.html"
}

// GetUsersLogs user logs
// @router /users/logs [get]
func (p *Plugin) GetUsersLogs() {
	p.LayoutDashboard()
	var items []Log
	_, err := orm.NewOrm().QueryTable(new(Log)).
		Filter("user_id", p.CurrentUser().ID).
		OrderBy("-updated_at").
		Limit(120).
		All(&items)
	p.Flash(nil, err)

	p.Data[TITLE] = Tr(p.Locale(), "nut.users.logs.title")
	p.Data["logs"] = items
	p.TplName = "nut/users/logs.html"
}

// DeleteUsersSignOut user sign out
// @router /users/sign-out [delete]
func (p *Plugin) DeleteUsersSignOut() {
	p.MustSignIn()
	p.DestroySession()
	AddLog(
		orm.NewOrm(),
		p.currentUser,
		p.Ctx.Input.Host(),
		Tr(p.locale, "logs.user.sign-out"),
	)
	p.Data["json"] = map[string]interface{}{"ok": true}
	p.ServeJSON()
}
