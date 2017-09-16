package forum

import (
	"net/http"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// IndexComments list all cards
// @router /comments [get]
func (p *Plugin) IndexComments() {
	p.LayoutDashboard()
	p.MustAdmin()
	var items []Comment
	if _, err := orm.NewOrm().QueryTable(new(Comment)).
		OrderBy("-updated_at").
		All(&items, "id", "body", "type", "user_id"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "forum.comments.index.title")

	p.TplName = "nut/comments/index.html"
}

// NewComment new card
// @router /comments/new [get]
func (p *Plugin) NewComment() {
	p.LayoutDashboard()

	var item Comment
	p.Data["item"] = item
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.new")
	p.Data["action"] = p.URLFor("nut.Plugin.CreateComment")
	p.TplName = "nut/comments/form.html"
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
		p.Redirect("nut.Plugin.IndexComments")
	} else {
		p.Redirect("nut.Plugin.NewComment")
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
	p.Data["action"] = p.URLFor("nut.Plugin.UpdateComment", ":id", id)
	p.Data["item"] = item
	p.TplName = "nut/comments/form.html"
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
		_, err = o.Update(&item, "body", "type")
	}

	if p.Flash(nil, err) {
		p.Redirect("nut.Plugin.IndexComment")
	} else {
		p.Redirect("nut.Plugin.EditComment", ":id", id)
	}
}

// DestroyComment remove
// @router /comments/:id [delete]
func (p *Plugin) DestroyComment() {
	p.MustSignIn()
	_, err := orm.NewOrm().QueryTable(new(Comment)).
		Filter("id", p.Ctx.Input.Param(":id")).
		Delete()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = nut.H{"ok": true}
	p.ServeJSON()
}
