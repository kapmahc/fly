package nut

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"
)

// Controller base controller
type Controller struct {
	beego.Controller
	Locale string
}

// Redirect http 302 redirect
func (p *Controller) Redirect(name string, args ...interface{}) {
	p.Controller.Redirect(p.URLFor(name, args...), http.StatusFound)
}

// Check write error flash if error
func (p *Controller) Check(e error) bool {
	if e == nil {
		return true
	}
	beego.Error(e)
	f := beego.NewFlash()
	f.Error(e.Error())
	f.Store(&p.Controller)
	return false
}

// Prepare prepare
func (p *Controller) Prepare() {
	beego.ReadFromRequest(&p.Controller)
	p.Data["xsrf"] = template.HTML(p.XSRFFormHTML())
	p.detectLocale()
}

// ParseForm parse form
func (p *Controller) ParseForm(fm interface{}) error {
	if er := p.Controller.ParseForm(fm); er != nil {
		return er
	}
	var va validation.Validation
	ok, er := va.Valid(fm)
	if er != nil {
		return er
	}
	if !ok {
		var msg []string
		for _, e := range va.Errors {
			msg = append(msg, fmt.Sprintf("%s: %s", e.Field, e.Message))
		}
		return errors.New(strings.Join(msg, "<br/>"))
	}
	return nil
}

// Abort http abort
func (p *Controller) Abort(s int, e error) {
	if e == nil {
		p.Controller.Abort("500")
	} else {
		beego.Error(e)
		p.CustomAbort(s, e.Error())
	}
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

// LayoutApplication use application layout
func (p *Controller) LayoutApplication() {
	// TODO
	p.Layout = "layouts/application/index.html"
}

// LayoutDashboard use dashboard layout
func (p *Controller) LayoutDashboard() {
	// TODO
	p.Layout = "layouts/dashboard/index.html"
}
