package forum

import (
	"net/http"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// IndexArticles list all cards
// @router /articles [get]
func (p *Plugin) IndexArticles() {
	p.LayoutDashboard()
	var items []Article
	if _, err := orm.NewOrm().QueryTable(new(Article)).
		OrderBy("-updated_at").
		All(&items, "id", "title"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "forum.articles.index.title")

	p.TplName = "nut/articles/index.html"
}

// NewArticle new card
// @router /articles/new [get]
func (p *Plugin) NewArticle() {
	p.LayoutDashboard()

	var item Article
	p.Data["item"] = item
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.new")
	p.Data["action"] = p.URLFor("nut.Plugin.CreateArticle")
	p.TplName = "nut/articles/form.html"
}

type fmArticle struct {
	Title string `form:"title" valid:"Required"`
	Body  string `form:"body" valid:"Required"`
	Type  string `form:"type" valid:"Required"`
}

// CreateArticle create
// @router /articles [post]
func (p *Plugin) CreateArticle() {
	p.MustSignIn()
	var fm fmArticle

	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().Insert(&Article{
			Title: fm.Title,
			Body:  fm.Body,
			Type:  fm.Type,
			User:  p.CurrentUser(),
		})
	}
	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexArticles")
	} else {
		p.Redirect("nut.Plugin.NewArticle")
	}
}

// EditArticle edit
// @router /articles/edit/:id [get]
func (p *Plugin) EditArticle() {
	p.LayoutDashboard()
	id := p.Ctx.Input.Param(":id")
	var item Article
	if err := orm.NewOrm().QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
		p.Abort(http.StatusForbidden, nil)
	}

	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.edit")
	p.Data["action"] = p.URLFor("nut.Plugin.UpdateArticle", ":id", id)
	p.Data["item"] = item
	p.TplName = "nut/articles/form.html"
}

// UpdateArticle update
// @router /articles/:id [post]
func (p *Plugin) UpdateArticle() {
	p.LayoutDashboard()
	var fm fmArticle

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	var item Article
	o := orm.NewOrm()
	if err == nil {
		err = o.QueryTable(new(Article)).
			Filter("id", id).
			One(&item)
	}
	if err == nil {
		if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
			err = nut.Te(p.Locale(), "errors.not-allow")
		}
	}
	if err == nil {
		_, err = o.Update(&item, "title", "body", "type")
	}

	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexArticle")
	} else {
		p.Redirect("nut.Plugin.EditArticle", ":id", id)
	}
}

// DestroyArticle remove
// @router /articles/:id [delete]
func (p *Plugin) DestroyArticle() {
	p.MustSignIn()
	_, err := orm.NewOrm().QueryTable(new(Article)).
		Filter("id", p.Ctx.Input.Param(":id")).
		Delete()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = nut.H{"ok": true}
	p.ServeJSON()
}
