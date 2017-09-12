package nut

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

// Controller base controller
type Controller struct {
	beego.Controller
	Locale string
}

// Prepare prepare
func (p *Controller) Prepare() {
	// detect lang
	p.detectLocale()
}

func (p *Controller) detectLocale() {
	write := false
	// 1. Check URL arguments.
	lang := p.Input().Get(LOCALE)

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = p.Ctx.GetCookie(LOCALE)
	} else {
		write = true
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		write = true
		al := p.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			lang = al[:5] // Only compare first 5 letters.
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "en-US"
	}

	// Save language information in cookies.
	if write {
		p.Ctx.SetCookie(LOCALE, lang, 1<<32-1, "/")
	}

	// Set language properties.
	p.Locale = lang
	p.Data[LOCALE] = lang
	p.Data["languages"] = i18n.ListLangs()

}

// ApplicationLayout application layout
func (p *Controller) ApplicationLayout() {
	// TODO
	p.Layout = "layouts/application/index.html"
}

// DashboardLayout dashboard layout
func (p *Controller) DashboardLayout() {
	// TODO
	p.Layout = "layouts/dashboard/index.html"
}
