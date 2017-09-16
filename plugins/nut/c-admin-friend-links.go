package nut

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
)

// IndexAdminFriendLinks list all friend-links
// @router /admin/friend-links [get]
func (p *Plugin) IndexAdminFriendLinks() {
	p.LayoutDashboard()
	p.MustAdmin()
	var items []FriendLink
	if _, err := orm.NewOrm().QueryTable(new(FriendLink)).
		OrderBy("sort_order").
		All(&items, "id", "title", "home", "logo", "sort_order"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.friend-links.index.title")

	p.TplName = "nut/admin/friend-links/index.html"
}

// NewAdminFriendLink new friend-link
// @router /admin/friend-links/new [get]
func (p *Plugin) NewAdminFriendLink() {
	p.LayoutDashboard()
	p.MustAdmin()
	var item FriendLink
	p.Data["item"] = item
	p.Data[TITLE] = Tr(p.Locale(), "buttons.new")
	p.SetSortOrders()
	p.Data["action"] = p.URLFor("nut.Plugin.CreateAdminFriendLink")
	p.TplName = "nut/admin/friend-links/form.html"
}

type fmFriendLink struct {
	Home      string `form:"home" valid:"Required"`
	Logo      string `form:"logo" valid:"Required"`
	Title     string `form:"title" valid:"Required"`
	SortOrder int    `form:"sortOrder"`
}

// CreateAdminFriendLink create
// @router /admin/friend-links [post]
func (p *Plugin) CreateAdminFriendLink() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmFriendLink

	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().Insert(&FriendLink{
			Home:      fm.Home,
			Logo:      fm.Logo,
			Title:     fm.Title,
			SortOrder: fm.SortOrder,
		})
	}
	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexAdminFriendLinks")
	} else {
		p.Redirect("nut.Plugin.NewAdminFriendLink")
	}
}

// EditAdminFriendLink edit
// @router /admin/friend-links/edit/:id [get]
func (p *Plugin) EditAdminFriendLink() {
	p.LayoutDashboard()
	p.MustAdmin()
	id := p.Ctx.Input.Param(":id")
	var item FriendLink
	if err := orm.NewOrm().QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data[TITLE] = Tr(p.Locale(), "buttons.edit")
	p.Data["action"] = p.URLFor("nut.Plugin.UpdateAdminFriendLink", ":id", id)
	p.Data["item"] = item
	p.SetSortOrders()
	p.TplName = "nut/admin/friend-links/form.html"
}

// UpdateAdminFriendLink update
// @router /admin/friend-links/:id [post]
func (p *Plugin) UpdateAdminFriendLink() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmFriendLink

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().QueryTable(new(FriendLink)).
			Filter("id", id).
			Update(orm.Params{
				"home":       fm.Home,
				"title":      fm.Title,
				"logo":       fm.Logo,
				"sort_order": fm.SortOrder,
				"updated_at": time.Now(),
			})
	}
	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexAdminFriendLink")
	} else {
		p.Redirect("nut.Plugin.EditAdminFriendLink", ":id", id)
	}
}

// DestroyAdminFriendLink remove
// @router /admin/friend-links/:id [delete]
func (p *Plugin) DestroyAdminFriendLink() {
	p.MustAdmin()
	_, err := orm.NewOrm().QueryTable(new(FriendLink)).
		Filter("id", p.Ctx.Input.Param(":id")).
		Delete()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = H{"ok": true}
	p.ServeJSON()
}
