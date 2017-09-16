package nut

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
)

// IndexAdminLinks list all links
// @router /admin/links [get]
func (p *Plugin) IndexAdminLinks() {
	p.LayoutDashboard()
	p.MustAdmin()
	var items []Link
	if _, err := orm.NewOrm().QueryTable(new(Link)).
		OrderBy("loc", "sort_order").
		All(&items, "id", "loc", "label", "href"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.links.index.title")
	p.TplName = "nut/admin/links/index.html"
}

// NewAdminLink new link
// @router /admin/links/new [get]
func (p *Plugin) NewAdminLink() {
	p.LayoutDashboard()
	p.MustAdmin()
	var item Link
	p.Data["item"] = item
	p.Data[TITLE] = Tr(p.Locale(), "buttons.new")
	p.TplName = "nut/admin/links/new.html"
}

type fmLink struct {
	Lable     string `form:"label" valid:"Required"`
	Loc       string `form:"loc" valid:"Required"`
	Href      string `form:"href" valid:"Required"`
	SortOrder int    `form:"sortOrder"`
}

// CreateAdminLink create
// @router /admin/links [post]
func (p *Plugin) CreateAdminLink() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmLink

	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().Insert(&Link{
			Label:     fm.Lable,
			Href:      fm.Href,
			Loc:       fm.Loc,
			SortOrder: fm.SortOrder,
		})
	}
	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexAdminLinks")
	} else {
		p.Redirect("nut.Plugin.NewAdminLink")
	}
}

// EditAdminLink edit
// @router /admin/links/edit/:id [get]
func (p *Plugin) EditAdminLink() {
	p.LayoutDashboard()
	p.MustAdmin()
	p.Data[TITLE] = Tr(p.Locale(), "buttons.edit")
	var item Link
	if err := orm.NewOrm().QueryTable(&item).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&item, "id", "loc", "href", "label", "sort_order"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	p.Data["item"] = item
	p.TplName = "nut/admin/links/edit.html"
}

// UpdateAdminLink update
// @router /admin/links/:id [post]
func (p *Plugin) UpdateAdminLink() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmLink

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().QueryTable(new(Link)).
			Filter("id", id).
			Update(orm.Params{
				"loc":        fm.Loc,
				"label":      fm.Lable,
				"href":       fm.Href,
				"sort_order": fm.SortOrder,
				"updated_at": time.Now(),
			})
	}
	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexAdminLink")
	} else {
		p.Redirect("nut.Plugin.EditAdminLink", ":id", id)
	}
}

// DestroyAdminLink remove
// @router /admin/links/:id [delete]
func (p *Plugin) DestroyAdminLink() {
	p.MustAdmin()
	_, err := orm.NewOrm().QueryTable(new(Link)).
		Filter("id", p.Ctx.Input.Param(":id")).
		Delete()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = H{"ok": true}
	p.ServeJSON()
}
