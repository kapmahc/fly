package nut

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"golang.org/x/text/language"
)

// GetUsersProfile users'profile
// @router /users/profile [get]
func (p *Plugin) GetUsersProfile() {
	p.LayoutDashboard()
	p.Data[TITLE] = Tr(p.Locale(), "nut.users.profile.title")
	p.TplName = "nut/users/profile.html"
}

type fmUserProfile struct {
	Name string `form:"name" valid:"Required"`
}

// PostUsersProfile users'profile
// @router /users/profile [post]
func (p *Plugin) PostUsersProfile() {
	p.LayoutDashboard()
	var fm fmUserProfile
	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().QueryTable(new(User)).Update(orm.Params{
			"updated_at": time.Now(),
			"name":       fm.Name,
		})
	}
	p.Flash(nil, err)
	p.Redirect("nut.Plugin.GetUsersProfile")
}

// GetUsersChangePassword change user password
// @router /users/change-password [get]
func (p *Plugin) GetUsersChangePassword() {
	p.LayoutDashboard()
	p.Data[TITLE] = Tr(p.Locale(), "nut.users.change-password.title")
	p.TplName = "nut/users/change-password.html"
}

type fmUserChangePassword struct {
	CurrentPassword      string `form:"currentPassword" valid:"Required"`
	NewPassword          string `form:"newPassword" valid:"MinSize(6)"`
	PasswordConfirmation string `form:"passwordConfirmation"`
}

func (p fmUserChangePassword) Valid(v *validation.Validation) {
	lang := language.AmericanEnglish.String()
	if p.NewPassword != p.PasswordConfirmation {
		v.SetError("PasswordConfirmation", Tr(lang, "nut.errors.user.passwords-not-match"))
	}
}

// PostUsersChangePassword change user's password
// @router /users/change-password [post]
func (p *Plugin) PostUsersChangePassword() {
	p.LayoutDashboard()
	var fm fmUserChangePassword
	err := p.ParseForm(&fm)
	if err == nil {
		if !HMAC().Chk([]byte(fm.CurrentPassword), []byte(p.CurrentUser().Password)) {
			err = Te(p.Locale(), "nut.errors.user.email-password-not-match")
		}
	}
	if err == nil {
		_, err = orm.NewOrm().QueryTable(new(User)).Update(orm.Params{
			"updated_at": time.Now(),
			"password":   string(HMAC().Sum([]byte(fm.NewPassword))),
		})
	}
	p.Flash(func() string {
		return Tr(p.Locale(), "helpers.success")
	}, err)
	p.Redirect("nut.Plugin.GetUsersChangePassword")
}

// GetUsersLogs user logs
// @router /users/logs [get]
func (p *Plugin) GetUsersLogs() {
	p.LayoutDashboard()
	var items []Log
	_, err := orm.NewOrm().QueryTable(new(Log)).
		Filter("user_id", p.CurrentUser().ID).
		OrderBy("-created_at").
		Limit(120).
		All(&items, "message", "created_at")
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
