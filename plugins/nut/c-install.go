package nut

import (
	"net/http"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"golang.org/x/text/language"
)

// GetInstall init database
// @router /install [get]
func (p *Plugin) GetInstall() {
	p.mustDbEmpty(orm.NewOrm())
	p.LayoutApplication()
	p.Data[TITLE] = Tr(p.Locale(), "nut.install.title")
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

func (p fmInstall) Valid(v *validation.Validation) {
	if p.Password != p.PasswordConfirmation {
		v.SetError("PasswordConfirmation", Tr(language.AmericanEnglish.String(), "nut.errors.user.passwords-not-match"))
	}
}

// PostInstall init database
// @router /install [post]
func (p *Plugin) PostInstall() {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	p.mustDbEmpty(o)

	lang := p.Locale()
	var fm fmInstall
	err := p.ParseForm(&fm)

	if err == nil {
		err = SetLocale(o, lang, "site.title", fm.Title)
	}
	if err == nil {
		err = SetLocale(o, lang, "site.subhead", fm.Subhead)
	}
	var user *User
	ip := p.Ctx.Input.IP()
	if err == nil {
		user, err = AddEmailUser(o, lang, ip, fm.Name, fm.Email, fm.Password)
	}
	if err == nil {
		err = confirmUser(o, lang, ip, user)
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
		o.Commit()
	} else {
		o.Rollback()
	}
	if p.Flash(nil, err) {
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
