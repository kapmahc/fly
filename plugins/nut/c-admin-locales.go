package nut

import (
	"net/http"

	"github.com/astaxie/beego/orm"
)

// IndexAdminLocales list all i18n items
// @router /admin/locales [get]
func (p *Plugin) IndexAdminLocales() {
	p.LayoutDashboard()
	p.MustAdmin()
	lang := p.Locale()
	var items []Locale
	if _, err := orm.NewOrm().QueryTable(new(Locale)).
		Filter("lang", lang).
		OrderBy("code").
		All(&items, "code", "message"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[TITLE] = Tr(lang, "nut.admin.locales.index.title")
	p.TplName = "nut/admin/locales/index.html"
}

// NewAdminLocale create a locale item
// @router /admin/locales/new [get]
func (p *Plugin) NewAdminLocale() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmLocale
	fm.Code = p.GetString("code", "")
	lang := p.Locale()
	if fm.Code == "" {
		p.Data[TITLE] = Tr(lang, "buttons.new")
	} else {
		fm.Message = Tr(lang, fm.Code)
		p.Data[TITLE] = Tr(lang, "buttons.edit")
	}
	p.Data["form"] = fm
	p.TplName = "nut/admin/locales/form.html"
}

type fmLocale struct {
	Code    string `form:"code" valid:"Required"`
	Message string `form:"message" valid:"Required"`
}

// CreateAdminLocale save a locale
// @router /admin/locales [post]
func (p *Plugin) CreateAdminLocale() {
	p.MustAdmin()
	var fm fmLocale

	lang := p.Locale()
	err := p.ParseForm(&fm)
	if err == nil {
		err = SetLocale(orm.NewOrm(), lang, fm.Code, fm.Message)
	}
	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexAdminLocales")
	} else {
		p.Redirect("nut.Plugin.NewAdminLocale", "code", fm.Code)
	}
}

// DestroyAdminLocale remove a locale
// @router /admin/locales/:code [delete]
func (p *Plugin) DestroyAdminLocale() {
	p.MustAdmin()
	_, err := orm.NewOrm().QueryTable(new(Locale)).
		Filter("code", p.Ctx.Input.Param(":code")).
		Filter("lang", p.Locale()).Delete()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = H{"ok": true}
	p.ServeJSON()
}
