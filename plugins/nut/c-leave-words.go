package nut

import (
	"github.com/astaxie/beego/orm"
)

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
