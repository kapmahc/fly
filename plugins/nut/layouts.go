package nut

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"
	"github.com/kapmahc/fly/plugins/nut/timeago"
)

// H hash
type H map[string]interface{}

// Controller base controller
type Controller struct {
	beego.Controller

	locale         string
	currentUser    *User
	isAdmin        bool
	dashboardMenus []H
}

// Locale get current locale
func (p *Controller) Locale() string {
	return p.locale
}

// CurrentUser get current user
func (p *Controller) CurrentUser() *User {
	return p.currentUser
}

// IsAdmin current user is admin?
func (p *Controller) IsAdmin() bool {
	return p.isAdmin
}

// Redirect http 302 redirect
func (p *Controller) Redirect(name string, args ...interface{}) {
	p.Controller.Redirect(p.URLFor(name, args...), http.StatusFound)
}

// HomeURL home url
func (p *Controller) HomeURL() string {
	req := p.Ctx.Request
	scheme := "http"
	if p.Ctx.Request.TLS != nil {
		scheme = scheme + "s"
	}
	return scheme + "://" + req.Host
}

// Flash write flash message
func (p *Controller) Flash(fn func() string, er error) bool {
	ok := false
	f := beego.NewFlash()
	if er == nil {
		if fn != nil {
			f.Notice(fn())
		}
		ok = true
	} else {
		beego.Error(er)
		f.Error(er.Error())
	}
	f.Store(&p.Controller)
	return ok
}

// Prepare prepare
func (p *Controller) Prepare() {
	beego.ReadFromRequest(&p.Controller)
	p.setXSRF()
	p.detectLocale()
	p.parseUserFromRequest()
}

func (p *Controller) setXSRF() {
	p.Data["xsrf_input"] = template.HTML(p.XSRFFormHTML())
	p.Data["xsrf_token"] = p.XSRFToken()
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
	const key = "locale"
	write := false
	// 1. Check URL arguments.
	lang := p.Input().Get(key)

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = p.Ctx.GetCookie(key)
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
		p.Ctx.SetCookie(key, lang, 1<<32-1, "/")
	}

	// Set language properties.
	p.locale = lang
	p.Data[key] = lang
	p.Data["languages"] = i18n.ListLangs()

}

// LayoutApplication use application layout
func (p *Controller) LayoutApplication() {
	// TODO
	p.Layout = "layouts/application/index.html"
}

func (p *Controller) parseUserFromRequest() {
	uid, ok := p.GetSession("uid").(string)
	if !ok {
		return
	}
	user, err := GetUserByUID(uid)
	if err != nil {
		return
	}

	if !user.IsConfirm() || user.IsLock() {
		return
	}
	p.currentUser = user
	p.isAdmin = Is(orm.NewOrm(), user.ID, RoleAdmin)
	p.Data["currentUser"] = user
	p.Data["isAdmin"] = p.isAdmin
}

// MustSignIn must-sign-in
func (p *Controller) MustSignIn() {
	if p.currentUser == nil {
		p.Abort(http.StatusForbidden, Te(p.locale, "nut.errors.user.please-sign-in"))
	}
}

// MustAdmin must has admin role
func (p *Controller) MustAdmin() {
	if !p.isAdmin {
		p.Abort(http.StatusForbidden, Te(p.locale, "errors.not-allow"))
	}
}

func init() {
	beego.AddFuncMap("timeago", timeago.FromTime)
}
