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

	locale  string
	user    *User
	isAdmin bool
}

// Locale get current locale
func (p *Controller) Locale() string {
	return p.locale
}

// Redirect http 302 redirect
func (p *Controller) Redirect(name string, args ...interface{}) {
	p.Controller.Redirect(p.URLFor(name, args...), http.StatusFound)
}

// HomeURL home url
func (p *Controller) HomeURL() string {
	req := p.Ctx.Request
	scheme := "https"
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

// LayoutDashboard use dashboard layout
func (p *Controller) LayoutDashboard() {
	// TODO
	p.Layout = "layouts/dashboard/index.html"
}

// func (p *Controller) parseUserFromRequest() {
// 	cm, err := JWT().ParseFromRequest(p.Ctx.Request)
// 	if err != nil {
// 		return
// 	}
// 	user, err := GetUserByUID(cm.Get(UID).(string))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !user.IsConfirm() {
// 		return nil, E(lng, "auth.errors.user.not-confirm")
// 	}
// 	if user.IsLock() {
// 		return nil, E(lng, "auth.errors.user.is-lock")
// 	}
// 	return user, nil
// }
//
// func (p *Controller) parseCurrentUser() error {
// 	if user, err := JWT().getUserFromRequest(c); err == nil {
// 		c.Set(CurrentUser, user)
// 		c.Set(IsAdmin, Is(user.ID, RoleAdmin))
// 	}
// 	return nil
// }
//
// // MustSignIn must-sign-in
// func (p *Controller) MustSignIn() error {
// 	if _, ok := c.MustGet(CurrentUser).(*User); ok {
// 		return nil
// 	}
// 	lng := c.MustGet(key).(string)
// 	return E(lng, "auth.errors.please-sign-in")
// }
//
// // MustAdminMiddleware must has admin role
// func (p *Controller) MustAdmin() error {
// 	if is, ok := c.MustGet(IsAdmin).(bool); ok && is {
// 		return nil
// 	}
// 	lng := c.MustGet(key).(string)
// 	return E(lng, "errors.not-allow")
// }
