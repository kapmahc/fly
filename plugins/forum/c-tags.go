package forum

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// IndexTags list all cards
// @router /tags [get]
func (p *Plugin) IndexTags() {
	p.LayoutDashboard()
	p.MustAdmin()
	var items []Tag
	if _, err := orm.NewOrm().QueryTable(new(Tag)).
		OrderBy("-updated_at").
		All(&items, "id", "name", "updated_at"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "forum.tags.index.title")

	p.TplName = "forum/tags/index.html"
}

// NewTag new card
// @router /tags/new [get]
func (p *Plugin) NewTag() {
	p.LayoutDashboard()
	p.MustAdmin()
	var item Tag
	p.Data["item"] = item
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.new")
	p.Data["action"] = p.URLFor("forum.Plugin.CreateTag")
	p.TplName = "forum/tags/form.html"
}

type fmTag struct {
	Name string `form:"name" valid:"Required"`
}

// CreateTag create
// @router /tags [post]
func (p *Plugin) CreateTag() {
	p.LayoutDashboard()
	p.MustAdmin()
	var fm fmTag

	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().Insert(&Tag{
			Name: fm.Name,
		})
	}

	if p.Flash(nil, err) {
		p.Redirect("forum.Plugin.IndexTags")
	} else {
		p.Redirect("forum.Plugin.NewTag")
	}
}

// ShowTag show
// @router /tags/:id [get]
func (p *Plugin) ShowTag() {
	p.LayoutApplication()
	id := p.Ctx.Input.Param(":id")
	var item Tag
	if err := orm.NewOrm().QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	p.Data[nut.TITLE] = item.Name
	p.Data["item"] = item
	p.TplName = "forum/tags/show.html"
}

// EditTag edit
// @router /tags/edit/:id [get]
func (p *Plugin) EditTag() {
	p.LayoutDashboard()
	p.MustAdmin()
	id := p.Ctx.Input.Param(":id")
	var item Tag
	if err := orm.NewOrm().QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.edit")
	p.Data["action"] = p.URLFor("forum.Plugin.UpdateTag", ":id", id)
	p.Data["item"] = item
	p.TplName = "forum/tags/form.html"
}

// UpdateTag update
// @router /tags/:id [post]
func (p *Plugin) UpdateTag() {
	p.LayoutDashboard()
	var fm fmTag

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().QueryTable(new(Tag)).
			Filter("id", id).
			Update(orm.Params{
				"name":       fm.Name,
				"updated_at": time.Now(),
			})
	}

	if p.Flash(nil, err) {
		p.Redirect("forum.Plugin.IndexTag")
	} else {
		p.Redirect("forum.Plugin.EditTag", ":id", id)
	}
}

// DestroyTag remove
// @router /tags/:id [delete]
func (p *Plugin) DestroyTag() {
	p.MustSignIn()
	_, err := orm.NewOrm().QueryTable(new(Tag)).
		Filter("id", p.Ctx.Input.Param(":id")).
		Delete()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = nut.H{"ok": true}
	p.ServeJSON()
}
