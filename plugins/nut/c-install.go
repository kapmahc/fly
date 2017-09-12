package nut

import (
	"errors"
	"net/http"

	"github.com/astaxie/beego/orm"
)

// GetInstall init database
// @router /install [get]
func (p *Plugin) GetInstall() {
	// p.mustDbEmpty(orm.NewOrm())
	p.LayoutApplication()
	p.Data[TITLE] = Tr(p.Locale, "nut.install.title")
	p.TplName = "nut/install.html"
}

type fmInstall struct {
	Title                string `form:"title" valid:"Required"`
	Subhead              string `form:"subhead" valid:"Required"`
	Name                 string `form:"name" valid:"Required"`
	Email                string `form:"email" valid:"Email"`
	Password             string `form:"password" valid:"MinSize(6)"`
	PasswordConfirmation string `form:"passwordConfirmation"`
}

// PostInstall init database
// @router /install [post]
func (p *Plugin) PostInstall() {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	p.mustDbEmpty(o)

	var fm fmInstall
	err := p.ParseForm(&fm)
	if err == nil {
		if fm.Password != fm.PasswordConfirmation {
			err = errors.New(Tr(p.Locale, "errors.passwords-not-match"))
		}
	}

	if err == nil {
		err = SetLocale(p.Locale, "site.title", fm.Title)
	}
	if err == nil {
		err = SetLocale(p.Locale, "site.subhead", fm.Subhead)
	}
	var user *User
	ip := p.Ctx.Input.IP()
	if err == nil {
		user, err = AddEmailUser(o, p.Locale, ip, fm.Name, fm.Email, fm.Password)
	}
	if err == nil {
		err = confirmUser(o, p.Locale, ip, user)
	}
	if err == nil {
		for _, r := range []string{RoleAdmin, RoleRoot} {
			var role *Role
			if role, err = GetRole(o, r, DefaultResourceType, DefaultResourceID); err != nil {
				break
			}
			if err = Allow(o, user, role, 100, 0, 0); err != nil {
				break
			}
		}
	}

	if err == nil {
		err = o.Commit()
	} else {
		err = o.Rollback()
	}
	if p.Check(err) {
		p.Redirect("nut.Plugin.GetHome")
	} else {
		p.Redirect("nut.Plugin.GetInstall")
	}
}

func (p *Plugin) mustDbEmpty(o orm.Ormer) {
	cnt, err := o.QueryTable(new(User)).Count()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	if cnt > 0 {
		p.Abort(http.StatusInternalServerError, nil)
	}
}
