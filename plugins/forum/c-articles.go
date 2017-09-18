package forum

import (
	"net/http"
	"strconv"
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
		Filter("user_id", p.CurrentUser().ID).
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
	p.setTagsData()
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.new")
	p.Data["action"] = p.URLFor("forum.Plugin.CreateArticle")
	p.TplName = "forum/articles/form.html"
}

func (p *Plugin) setTagsData() {
	var items []Tag
	if _, err := orm.NewOrm().QueryTable(new(Tag)).
		All(&items, "id", "name"); err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["tags"] = items
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
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	var fm fmArticle

	err := p.ParseForm(&fm)
	var item Article
	if err == nil {
		item = Article{
			Title: fm.Title,
			Body:  fm.Body,
			Type:  fm.Type,
			User:  p.CurrentUser(),
		}
		_, err = o.Insert(&item)
	}
	if err == nil {
		var tags []interface{}
		for _, t := range p.Ctx.Request.Form["tags"] {
			var id int
			id, err = strconv.Atoi(t)
			if err != nil {
				break
			}
			tags = append(tags, &Tag{ID: uint(id)})
		}
		if err == nil {
			_, err = o.QueryM2M(&item, "Tags").Add(tags...)
		}
	}
	if err == nil {
		o.Commit()
	} else {
		o.Rollback()
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
	o := orm.NewOrm()
	if err := o.QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if _, err := o.LoadRelated(&item, "Comments"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if _, err := o.LoadRelated(&item, "Tags"); err != nil {
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
	o := orm.NewOrm()
	if err := o.QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if _, err := o.LoadRelated(&item, "Tags"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
		p.Abort(http.StatusForbidden, nil)
	}

	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.edit")
	p.Data["action"] = p.URLFor("forum.Plugin.UpdateArticle", ":id", id)
	p.Data["item"] = item
	p.setTagsData()
	p.TplName = "forum/articles/form.html"
}

// UpdateArticle update
// @router /articles/:id [post]
func (p *Plugin) UpdateArticle() {
	p.LayoutDashboard()
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	var fm fmArticle

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	var item Article
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

	if err == nil {
		_, err = o.QueryM2M(&item, "Tags").Clear()
	}

	if err == nil {
		var tags []interface{}
		for _, t := range p.Ctx.Request.Form["tags"] {
			var id int
			id, err = strconv.Atoi(t)
			if err != nil {
				break
			}
			tags = append(tags, &Tag{ID: uint(id)})
		}
		if err == nil {
			_, err = o.QueryM2M(&item, "Tags").Add(tags...)
		}
	}

	if err == nil {
		o.Commit()
	} else {
		o.Rollback()
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
	o := orm.NewOrm()
	var item Article
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	err := o.QueryTable(&item).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&item)
	if err == nil {
		if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
			err = nut.Te(p.Locale(), "errors.not-allow")
		}
	}
	if err == nil {
		_, err = o.QueryTable(new(Comment)).Filter("article_id", item.ID).Delete()
	}
	if err == nil {
		_, err = o.QueryM2M(&item, "Tags").Clear()
	}
	if err == nil {
		o.Delete(&item)
	}

	if err == nil {
		o.Commit()
	} else {
		o.Rollback()
	}
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = nut.H{"ok": true}
	p.ServeJSON()
}
