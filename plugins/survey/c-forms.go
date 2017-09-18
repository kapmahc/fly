package survey

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/google/uuid"
	"github.com/kapmahc/fly/plugins/nut"
)

// IndexForms list all cards
// @router /forms [get]
func (p *Plugin) IndexForms() {
	p.LayoutDashboard()
	var items []Form
	if _, err := orm.NewOrm().QueryTable(new(Form)).
		OrderBy("-updated_at").
		All(&items, "id", "title", "updated_at"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "survey.forms.index.title")

	p.TplName = "survey/forms/index.html"
}

// NewForm new card
// @router /forms/new [get]
func (p *Plugin) NewForm() {
	p.LayoutDashboard()

	var item Form
	p.Data["item"] = item
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.new")
	p.Data["action"] = p.URLFor("survey.Plugin.CreateForm")
	p.TplName = "survey/forms/form.html"
}

type fmForm struct {
	Title string `form:"title" valid:"Required"`
	Body  string `form:"body" valid:"Required"`
	Type  string `form:"type" valid:"Required"`
}

// CreateForm create
// @router /forms [post]
func (p *Plugin) CreateForm() {
	p.MustSignIn()
	var fm fmForm

	err := p.ParseForm(&fm)
	var item Form
	if err == nil {
		item = Form{
			Title: fm.Title,
			Body:  fm.Body,
			Type:  fm.Type,
			UID:   uuid.New().String(),
			Mode:  "public",
			User:  p.CurrentUser(),
		}
		_, err = orm.NewOrm().Insert(&item)
	}

	if p.Flash(nil, err) {
		p.Redirect("survey.Plugin.IndexForms")
	} else {
		p.Redirect("survey.Plugin.NewForm")
	}
}

// ShowForm show
// @router /forms/:id [get]
func (p *Plugin) ShowForm() {
	p.LayoutApplication()
	id := p.Ctx.Input.Param(":id")
	var item Form
	o := orm.NewOrm()
	if err := o.QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if _, err := o.LoadRelated(&item, "Fields"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	p.Data[nut.TITLE] = item.Title
	p.Data["item"] = item
	p.TplName = "survey/forms/show.html"
}

// EditForm edit
// @router /forms/edit/:id [get]
func (p *Plugin) EditForm() {
	p.LayoutDashboard()
	id := p.Ctx.Input.Param(":id")
	var item Form
	o := orm.NewOrm()
	if err := o.QueryTable(&item).
		Filter("id", id).
		One(&item); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
		p.Abort(http.StatusForbidden, nil)
	}

	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.edit")
	p.Data["action"] = p.URLFor("survey.Plugin.UpdateForm", ":id", id)
	p.Data["item"] = item
	p.TplName = "survey/forms/form.html"
}

// UpdateForm update
// @router /forms/:id [post]
func (p *Plugin) UpdateForm() {
	p.LayoutDashboard()
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	var fm fmForm

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	var item Form
	if err == nil {
		err = o.QueryTable(new(Form)).
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
		o.Commit()
	} else {
		o.Rollback()
	}

	if p.Flash(nil, err) {
		p.Redirect("survey.Plugin.IndexForm")
	} else {
		p.Redirect("survey.Plugin.EditForm", ":id", id)
	}
}

// DestroyForm remove
// @router /forms/:id [delete]
func (p *Plugin) DestroyForm() {
	p.MustSignIn()
	o := orm.NewOrm()
	var item Form
	if err := o.Begin(); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	err := o.QueryTable(&item).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&item)
	if err == nil {
		_, err = o.QueryTable(new(Field)).Filter("form_id", item.ID).Delete()
	}
	if err == nil {
		_, err = o.QueryTable(new(Record)).Filter("form_id", item.ID).Delete()
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
