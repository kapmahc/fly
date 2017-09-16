package nut

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
)

// IndexAdminCards list all cards
// @router /admin/cards [get]
func (p *Plugin) IndexAdminCards() {
	p.LayoutDashboard()
	p.MustAdmin()
	var items []Card
	if _, err := orm.NewOrm().QueryTable(new(Card)).
		OrderBy("loc", "sort_order").
		All(&items, "id", "loc", "title", "href", "sort_order"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.cards.index.title")

	p.TplName = "nut/admin/cards/index.html"
}

// NewAdminCard new card
// @router /admin/cards/new [get]
func (p *Plugin) NewAdminCard() {
	p.LayoutDashboard()
	p.MustAdmin()
	var item Card
	p.Data["item"] = item
	p.Data[TITLE] = Tr(p.Locale(), "buttons.new")
	p.SetSortOrders()
	p.Data["action"] = p.URLFor("nut.Plugin.CreateAdminCard")
	p.TplName = "nut/admin/cards/form.html"
}

type fmCard struct {
	Title     string `form:"title" valid:"Required"`
	Summary   string `form:"summary" valid:"Required"`
	Type      string `form:"type" valid:"Required"`
	Action    string `form:"action" valid:"Required"`
	Logo      string `form:"logo" valid:"Required"`
	Loc       string `form:"loc" valid:"Required"`
	Href      string `form:"href" valid:"Required"`
	SortOrder int    `form:"sortOrder"`
}

// CreateAdminCard create
// @router /admin/cards [post]
func (p *Plugin) CreateAdminCard() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmCard

	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().Insert(&Card{
			Title:     fm.Title,
			Summary:   fm.Summary,
			Type:      fm.Type,
			Action:    fm.Action,
			Logo:      fm.Logo,
			Href:      fm.Href,
			Loc:       fm.Loc,
			SortOrder: fm.SortOrder,
		})
	}
	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexAdminCards")
	} else {
		p.Redirect("nut.Plugin.NewAdminCard")
	}
}

// EditAdminCard edit
// @router /admin/cards/edit/:id [get]
func (p *Plugin) EditAdminCard() {
	p.LayoutDashboard()
	p.MustAdmin()
	id := p.Ctx.Input.Param(":id")
	var item Card
	if err := orm.NewOrm().QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data[TITLE] = Tr(p.Locale(), "buttons.edit")
	p.Data["action"] = p.URLFor("nut.Plugin.UpdateAdminCard", ":id", id)
	p.Data["item"] = item
	p.SetSortOrders()
	p.TplName = "nut/admin/cards/form.html"
}

// UpdateAdminCard update
// @router /admin/cards/:id [post]
func (p *Plugin) UpdateAdminCard() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmCard

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().QueryTable(new(Card)).
			Filter("id", id).
			Update(orm.Params{
				"title":      fm.Title,
				"summary":    fm.Summary,
				"type":       fm.Type,
				"action":     fm.Action,
				"logo":       fm.Logo,
				"loc":        fm.Loc,
				"href":       fm.Href,
				"sort_order": fm.SortOrder,
				"updated_at": time.Now(),
			})
	}
	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexAdminCard")
	} else {
		p.Redirect("nut.Plugin.EditAdminCard", ":id", id)
	}
}

// DestroyAdminCard remove
// @router /admin/cards/:id [delete]
func (p *Plugin) DestroyAdminCard() {
	p.MustAdmin()
	_, err := orm.NewOrm().QueryTable(new(Card)).
		Filter("id", p.Ctx.Input.Param(":id")).
		Delete()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = H{"ok": true}
	p.ServeJSON()
}
