package nut

import "github.com/astaxie/beego/orm"

// IndexAdminUsers list all users
// @router /admin/users [get]
func (p *Plugin) IndexAdminUsers() {
	p.LayoutDashboard()
	p.MustAdmin()
	var items []User
	_, err := orm.NewOrm().QueryTable(new(User)).
		OrderBy("-current_sign_in_at").
		All(&items, "id", "email", "name", "last_sign_in_at", "last_sign_in_ip", "current_sign_in_ip", "current_sign_in_at")
	p.Flash(nil, err)

	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.users.index.title")
	p.Data["users"] = items
	p.TplName = "nut/admin/users/index.html"
}
