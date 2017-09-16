package nut

import (
	"net/http"

	"github.com/astaxie/beego/orm"
)

// IndexLeaveWords  index leave-words
// @router /leave-words [get]
func (p *Plugin) IndexLeaveWords() {
	p.LayoutDashboard()
	p.MustAdmin()
	var items []LeaveWord
	if _, err := orm.NewOrm().QueryTable(new(LeaveWord)).
		OrderBy("-created_at").
		All(&items); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	p.Data["items"] = items
	p.Data[TITLE] = Tr(p.Locale(), "nut.leave-words.index.title")
	p.TplName = "nut/leave-words/index.html"
}

// NewLeaveWord new leave word
// @router /leave-words/new [get]
func (p *Plugin) NewLeaveWord() {
	p.LayoutApplication()
	p.Data[TITLE] = Tr(p.Locale(), "nut.leave-words.new.title")
	p.TplName = "nut/leave-words/new.html"
}

type fmLeaveWord struct {
	Body string `form:"body" valid:"Required"`
	Type string `form:"type" valid:"Required"`
}

// CreateLeaveWord create leave-word
// @router /leave-words [post]
func (p *Plugin) CreateLeaveWord() {
	var fm fmLeaveWord
	err := p.ParseForm(&fm)
	if err == nil {
		_, err = orm.NewOrm().Insert(&LeaveWord{
			Body: fm.Body,
			Type: fm.Type,
		})
	}

	p.Flash(func() string {
		return Tr(p.Locale(), "helpers.success")
	}, err)
	p.Redirect("nut.Plugin.NewLeaveWord")
}

// DestroyLeaveWord remove
// @router /leave-words/:id [delete]
func (p *Plugin) DestroyLeaveWord() {
	p.MustAdmin()
	_, err := orm.NewOrm().QueryTable(new(LeaveWord)).
		Filter("id", p.Ctx.Input.Param(":id")).
		Delete()
	if err != nil {
		p.Abort(http.StatusOK, err)
	}
	p.Data["json"] = H{"ok": true}
	p.ServeJSON()
}
