package nut

import (
	"errors"
	"net/http"

	"github.com/astaxie/beego/orm"
)

// GetInstall init database
// @router /install [get]
func (p *Plugin) GetInstall() {
	p.mustDbEmpty()
	p.LayoutApplication()
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
	p.mustDbEmpty()

	var fm fmInstall
	err := p.ParseForm(&fm)
	if err == nil {
		if fm.Password != fm.PasswordConfirmation {
			err = errors.New(Tr(p.Locale, "errors.passwords-not-match"))
		}
	}

	// TODO
	p.Check(err)
	p.Redirect("/install", http.StatusFound)
}

func (p *Plugin) mustDbEmpty() {
	cnt, err := orm.NewOrm().QueryTable(new(User)).Count()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	if cnt > 0 {
		p.Abort(http.StatusInternalServerError, nil)
	}
}
