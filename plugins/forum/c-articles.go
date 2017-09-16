package forum

import (
	"net/http"
	"time"

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
		All(&items, "id", "title", "updated_at"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "forum.articles.index.title")

	p.TplName = "forum/articles/index.html"
}

// NewArticle new card
// @router /articles/new [get]
func (p *Plugin) NewArticle() {
	p.LayoutDashboard()

	var item Article
	p.Data["item"] = item
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.new")
	p.Data["action"] = p.URLFor("forum.Plugin.CreateArticle")
	p.TplName = "forum/articles/form.html"
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
		p.Redirect("forum.Plugin.IndexArticles")
	} else {
		p.Redirect("forum.Plugin.NewArticle")
	}
}

// ShowArticle show
// @router /articles/:id [get]
func (p *Plugin) ShowArticle() {
	p.LayoutApplication()
	id := p.Ctx.Input.Param(":id")
	var item Article
	if err := orm.NewOrm().QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	p.Data[nut.TITLE] = item.Title
	p.Data["item"] = item
	p.TplName = "forum/articles/show.html"
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
	p.Data["action"] = p.URLFor("forum.Plugin.UpdateArticle", ":id", id)
	p.Data["item"] = item
	p.TplName = "forum/articles/form.html"
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
		item.Title = fm.Title
		item.Body = fm.Body
		item.Type = fm.Type
		item.UpdatedAt = time.Now()
		_, err = o.Update(&item, "title", "body", "type", "updated_at")
	}

	if p.Flash(nil, err) {
		p.Redirect("forum.Plugin.IndexArticle")
	} else {
		p.Redirect("forum.Plugin.EditArticle", ":id", id)
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
