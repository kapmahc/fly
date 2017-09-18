package forum

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// IndexComments list all cards
// @router /comments [get]
func (p *Plugin) IndexComments() {
	p.LayoutDashboard()
	var items []Comment
	if _, err := orm.NewOrm().QueryTable(new(Comment)).
		OrderBy("-updated_at").
		Filter("user_id", p.CurrentUser().ID).
		All(&items, "id", "article_id", "body", "type", "user_id", "updated_at"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "forum.comments.index.title")

	p.TplName = "forum/comments/index.html"
}

// NewComment new card
// @router /comments/new [get]
func (p *Plugin) NewComment() {
	p.LayoutDashboard()
	id, err := p.GetInt("articleId", 0)
	if err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	var item Comment
	item.Article = &Article{ID: uint(id)}
	p.Data["item"] = item
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.new")
	p.Data["action"] = p.URLFor("forum.Plugin.CreateComment")
	p.TplName = "forum/comments/form.html"
}

type fmComment struct {
	ArticleID uint   `form:"articleId" valid:"Required"`
	Body      string `form:"body" valid:"Required"`
	Type      string `form:"type" valid:"Required"`
}

// CreateComment create
// @router /comments [post]
func (p *Plugin) CreateComment() {
	p.MustSignIn()
	var fm fmComment

	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().Insert(&Comment{
			Article: &Article{ID: fm.ArticleID},
			Body:    fm.Body,
			Type:    fm.Type,
			User:    p.CurrentUser(),
		})
	}
	if p.Flash(nil, err) {
		p.Redirect("forum.Plugin.ShowArticle", ":id", fm.ArticleID)
	} else {
		p.Redirect("forum.Plugin.NewComment")
	}
}

// EditComment edit
// @router /comments/edit/:id [get]
func (p *Plugin) EditComment() {
	p.LayoutDashboard()
	id := p.Ctx.Input.Param(":id")
	var item Comment
	if err := orm.NewOrm().QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
		p.Abort(http.StatusForbidden, nil)
	}

	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.edit")
	p.Data["action"] = p.URLFor("forum.Plugin.UpdateComment", ":id", id)
	p.Data["item"] = item
	p.TplName = "forum/comments/form.html"
}

// UpdateComment update
// @router /comments/:id [post]
func (p *Plugin) UpdateComment() {
	p.MustSignIn()
	var fm fmComment

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	var item Comment
	o := orm.NewOrm()
	if err == nil {
		err = o.QueryTable(new(Comment)).
			Filter("id", id).
			One(&item)
	}
	if err == nil {
		if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
			err = nut.Te(p.Locale(), "errors.not-allow")
		}
	}
	if err == nil {
		item.Body = fm.Body
		item.Type = fm.Type
		item.UpdatedAt = time.Now()
		_, err = o.Update(&item, "body", "type")
	}

	if p.Flash(nil, err) {
		p.Redirect("forum.Plugin.ShowArticle", ":id", fm.ArticleID)
	} else {
		p.Redirect("forum.Plugin.EditComment", ":id", id)
	}
}

// DestroyComment remove
// @router /comments/:id [delete]
func (p *Plugin) DestroyComment() {
	p.MustSignIn()
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	var item Comment
	err := o.QueryTable(&item).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&item)
	if err == nil {
		if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
			err = nut.Te(p.Locale(), "errors.not-allow")
		}
	}
	if err == nil {
		_, err = o.Delete(&item)
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
