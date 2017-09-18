package survey

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

func (p *Plugin) getField() (*Field, error) {
	o := orm.NewOrm()
	var item Field
	if err := o.QueryTable(&item).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&item); err != nil {
		return nil, err
	}
	if _, err := o.LoadRelated(&item, "Form"); err != nil {
		return nil, err
	}
	if item.Form.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
		return nil, nut.Te(p.Locale(), "errors.not-allow")
	}

	return &item, nil
}

func (p *Plugin) getForm() (*Form, error) {
	fid, err := p.GetInt("formId")
	if err != nil {
		return nil, err
	}
	var item Form
	o := orm.NewOrm()
	if err := o.QueryTable(&item).
		Filter("id", fid).
		One(&item); err != nil {
		return nil, err
	}
	if item.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
		return nil, nut.Te(p.Locale(), "errors.not-allow")
	}

	return &item, nil

}

// IndexFields list all cards
// @router /fields [get]
func (p *Plugin) IndexFields() {
	p.LayoutDashboard()

	form, err := p.getForm()
	if err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	var items []Field
	if _, err := orm.NewOrm().QueryTable(new(Field)).
		OrderBy("sort_order").
		Filter("form_id", form.ID).
		All(&items); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data["form"] = form
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "survey.fields.index.title")

	p.TplName = "survey/fields/index.html"
}

// NewField new card
// @router /fields/new [get]
func (p *Plugin) NewField() {
	p.LayoutDashboard()
	form, err := p.getForm()
	if err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	var item Field
	item.Form = form
	p.Data["item"] = item
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.new")
	p.Data["action"] = p.URLFor("survey.Plugin.CreateField", "formId", form.ID)
	p.SetSortOrders()
	p.setFieldTypes()
	p.TplName = "survey/fields/form.html"
}

type fmField struct {
	Label     string `form:"label" valid:"Required"`
	Name      string `form:"name" valid:"Required"`
	Body      string `form:"body"`
	Type      string `form:"type" valid:"Required"`
	Value     string `form:"value"`
	Required  bool   `form:"required"`
	SortOrder int    `form:"sortOrder"`
}

// CreateField create
// @router /fields [post]
func (p *Plugin) CreateField() {
	p.MustSignIn()

	var fm fmField
	err := p.ParseForm(&fm)
	var item Field
	var form *Form
	if err == nil {
		form, err = p.getForm()
	}
	if err == nil {
		item = Field{
			Name:      fm.Name,
			Label:     fm.Label,
			Body:      fm.Body,
			Type:      fm.Type,
			Value:     fm.Value,
			Required:  fm.Required,
			SortOrder: fm.SortOrder,
			Form:      &Form{ID: form.ID},
		}
		_, err = orm.NewOrm().Insert(&item)
	}

	if p.Flash(nil, err) {
		p.Redirect("survey.Plugin.IndexFields", "formId", p.GetString("formId"))
	} else {
		p.Redirect("survey.Plugin.NewField", "formId", p.GetString("formId"))
	}
}

// EditField edit
// @router /fields/edit/:id [get]
func (p *Plugin) EditField() {
	p.LayoutDashboard()
	item, err := p.getField()
	if err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "buttons.edit")
	p.Data["action"] = p.URLFor("survey.Plugin.UpdateField", ":id", item.ID)
	p.Data["item"] = item
	p.SetSortOrders()
	p.setFieldTypes()
	p.TplName = "survey/fields/form.html"
}

func (p *Plugin) setFieldTypes() {
	p.Data["types"] = []string{
		"text",
		"textarea",
		"select",
		"checkboxs",
	}
}

// UpdateField update
// @router /fields/:id [post]
func (p *Plugin) UpdateField() {
	p.LayoutDashboard()

	var fm fmField

	id := p.Ctx.Input.Param(":id")
	err := p.ParseForm(&fm)
	var item *Field
	if err == nil {
		item, err = p.getField()
	}

	if err == nil {
		item.Name = fm.Name
		item.Label = fm.Label
		item.Body = fm.Body
		item.Type = fm.Type
		item.Value = fm.Value
		item.Required = fm.Required
		item.SortOrder = fm.SortOrder
		item.UpdatedAt = time.Now()
		_, err = orm.NewOrm().Update(item,
			"name", "label", "body", "type", "value", "required", "sort_order",
			"updated_at")
	}

	if p.Flash(nil, err) {
		p.Redirect("survey.Plugin.IndexFields", "formId", item.Form.ID)
	} else {
		p.Redirect("survey.Plugin.EditField", ":id", id)
	}
}

// DestroyField remove
// @router /fields/:id [delete]
func (p *Plugin) DestroyField() {
	p.MustSignIn()

	item, err := p.getField()
	if err == nil {
		orm.NewOrm().Delete(item)
	}

	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = nut.H{"ok": true}
	p.ServeJSON()
}
