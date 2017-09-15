package nut

import (
	"net/http"

	"golang.org/x/text/language"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

// GetAdminSiteInfo edit site info
// @router /admin/site/info [get]
func (p *Plugin) GetAdminSiteInfo() {
	p.LayoutDashboard()
	p.MustAdmin()
	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.site.info.title")
	p.TplName = "nut/admin/site/info.html"
}

type fmSiteInfo struct {
	Title       string `form:"title" valid:"Required"`
	Subhead     string `form:"subhead" valid:"Required"`
	Keywords    string `form:"keywords" valid:"Required"`
	Description string `form:"description" valid:"Required"`
	Copyright   string `form:"copyright" valid:"Required"`
}

// PostAdminSiteInfo update site info
// @router /admin/site/info [post]
func (p *Plugin) PostAdminSiteInfo() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmSiteInfo
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	lang := p.Locale()
	err := p.ParseForm(&fm)
	if err == nil {
		for k, v := range map[string]string{
			"title":       fm.Title,
			"subhead":     fm.Subhead,
			"keywords":    fm.Keywords,
			"description": fm.Description,
			"copyright":   fm.Copyright,
		} {
			if err = SetLocale(o, lang, "site."+k, v); err != nil {
				break
			}
		}
	}
	if err == nil {
		o.Commit()
	} else {
		o.Rollback()
	}
	p.Flash(nil, err)
	p.Redirect("nut.Plugin.GetAdminSiteInfo")
}

// GetAdminSiteAuthor edit site author
// @router /admin/site/author [get]
func (p *Plugin) GetAdminSiteAuthor() {
	p.LayoutDashboard()
	p.MustAdmin()
	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.site.author.title")
	p.TplName = "nut/admin/site/author.html"
}

type fmSiteAuthor struct {
	Name  string `form:"name" valid:"Required"`
	Email string `form:"email" valid:"Required"`
}

// PostAdminSiteAuthor update author info
// @router /admin/site/author [post]
func (p *Plugin) PostAdminSiteAuthor() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmSiteAuthor
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	lang := p.Locale()
	err := p.ParseForm(&fm)
	if err == nil {
		for k, v := range map[string]string{
			"name":  fm.Name,
			"email": fm.Email,
		} {
			if err = SetLocale(o, lang, "site.author."+k, v); err != nil {
				break
			}
		}
	}
	if err == nil {
		o.Commit()
	} else {
		o.Rollback()
	}
	p.Flash(nil, err)
	p.Redirect("nut.Plugin.GetAdminSiteAuthor")
}

// GetAdminSiteSeo edit site seo
// @router /admin/site/seo [get]
func (p *Plugin) GetAdminSiteSeo() {
	p.LayoutDashboard()
	p.MustAdmin()

	var fm fmSiteSeo
	Get("site.google.verify-code", &fm.GoogleVerifyCode)
	Get("site.baidu.verify-code", &fm.BaiduVerifyCode)
	p.Data["form"] = fm

	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.site.seo.title")
	p.TplName = "nut/admin/site/seo.html"
}

type fmSiteSeo struct {
	GoogleVerifyCode string `form:"googleVerifyCode"`
	BaiduVerifyCode  string `form:"baiduVerifyCode"`
}

// PostAdminSiteSeo update author seo
// @router /admin/site/seo [post]
func (p *Plugin) PostAdminSiteSeo() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmSiteSeo
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	err := p.ParseForm(&fm)
	if err == nil {
		for k, v := range map[string]string{
			"google.verify-code": fm.GoogleVerifyCode,
			"baidu.verify-code":  fm.BaiduVerifyCode,
		} {
			if err = Set(o, "site."+k, v, false); err != nil {
				break
			}
		}
	}
	if err == nil {
		o.Commit()
	} else {
		o.Rollback()
	}
	p.Flash(nil, err)
	p.Redirect("nut.Plugin.GetAdminSiteSeo")
}

// GetAdminSiteSMTP edit site smtp
// @router /admin/site/smtp [get]
func (p *Plugin) GetAdminSiteSMTP() {
	p.LayoutDashboard()
	p.MustAdmin()
	smtp := map[string]interface{}{}
	if err := Get("site.smtp", &smtp); err != nil {
		smtp["host"] = "localhost"
		smtp["port"] = 25
		smtp["sender"] = "who-am-i@change-me.com"
	} else {
		smtp["password"] = ""
		smtp["port"] = int(smtp["port"].(float64))
	}

	p.Data["ports"] = []int{25, 465, 587}
	p.Data["form"] = smtp
	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.site.smtp.title")
	p.TplName = "nut/admin/site/smtp.html"
}

type fmSiteSMTP struct {
	Host                 string `form:"host" valid:"Required"`
	Port                 int    `form:"port"`
	Sender               string `form:"sender" valid:"Email"`
	Password             string `form:"password" valid:"Required"`
	PasswordConfirmation string `form:"passwordConfirmation" valid:"Required"`
}

func (p fmSiteSMTP) Valid(v *validation.Validation) {
	if p.Password != p.PasswordConfirmation {
		v.SetError("PasswordConfirmation", Tr(language.AmericanEnglish.String(), "nut.errors.user.passwords-not-match"))
	}
}

// PostAdminSiteSMTP update author smtp
// @router /admin/site/smtp [post]
func (p *Plugin) PostAdminSiteSMTP() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmSiteSMTP
	err := p.ParseForm(&fm)
	if err == nil {
		err = Set(
			orm.NewOrm(),
			"site.smtp",
			H{
				"host":     fm.Host,
				"port":     fm.Port,
				"sender":   fm.Sender,
				"password": fm.Password,
			},
			true,
		)
	}
	p.Flash(nil, err)
	p.Redirect("nut.Plugin.GetAdminSiteSMTP")
}
