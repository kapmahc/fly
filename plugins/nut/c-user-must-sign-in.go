package nut

import "github.com/astaxie/beego/orm"

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
