package survey

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// IndexRecords index
// @router /records [get]
func (p *Plugin) IndexRecords() {
	p.LayoutDashboard()

	form, err := p.getForm()
	if err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	o := orm.NewOrm()
	var fields []Field
	if _, err := o.QueryTable(new(Field)).
		OrderBy("sort_order").
		Filter("form_id", form.ID).
		All(&fields, "name"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	var records []Record
	if _, err := orm.NewOrm().QueryTable(new(Record)).
		OrderBy("-updated_at").
		Filter("form_id", form.ID).
		All(&records, "id", "updated_at", "value"); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	var values []nut.H
	for _, r := range records {
		v := map[string]interface{}{}
		if err := json.Unmarshal([]byte(r.Value), &v); err != nil {
			beego.Error(err)
			break
		}
		v["updatedAt"] = r.UpdatedAt.Format(time.RFC822)
		v["id"] = r.ID
		values = append(values, v)
	}

	p.Data["form"] = form
	p.Data["fields"] = fields
	p.Data["items"] = values
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "survey.records.index.title")

	p.TplName = "survey/records/index.html"
}

// CreateRecord create
// @router /records [post]
func (p *Plugin) CreateRecord() {
	fid, err := p.GetInt("formId")
	if err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	var form Form
	o := orm.NewOrm()
	if err == nil {
		err = o.QueryTable(&form).
			Filter("id", fid).
			One(&form)
	}
	if err == nil {
		_, err = o.LoadRelated(&form, "Fields")
	}
	values := nut.H{}
	for _, f := range form.Fields {
		val := p.Ctx.Request.Form[f.Name]
		if f.Required && len(val) == 0 {
			err = errors.New(f.Name + " mustn't empty")
			break
		}
		values[f.Name] = val
	}
	var val []byte
	if err == nil {
		val, err = json.Marshal(values)
	}
	if err == nil {
		_, err = o.Insert(&Record{
			Form:  &form,
			Value: string(val),
		})
	}

	p.Flash(func() string {
		return nut.Tr(p.Locale(), "helper.success")
	}, err)
	p.Redirect("survey.Plugin.ShowForm", ":id", fid)
}

// DestroyRecord remove
// @router /records/:id [delete]
func (p *Plugin) DestroyRecord() {
	p.MustSignIn()

	o := orm.NewOrm()
	var item Record
	err := o.QueryTable(&item).
		Filter("id", p.Ctx.Input.Param(":id")).
		One(&item)
	if err == nil {
		_, err = o.LoadRelated(&item, "Form")
	}
	if err == nil {
		if item.Form.User.ID != p.CurrentUser().ID && !p.IsAdmin() {
			err = nut.Te(p.Locale(), "errors.not-allow")
		}
	}

	if err == nil {
		orm.NewOrm().Delete(&item)
	}

	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = nut.H{"ok": true}
	p.ServeJSON()
}
