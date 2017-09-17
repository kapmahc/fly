package nut

import (
	"net/http"

	"github.com/astaxie/beego/orm"
)

// IndexAttachments list
// @router /attachments [get]
func (p *Plugin) IndexAttachments() {
	p.LayoutDashboard()
	var items []Attachment
	if _, err := orm.NewOrm().QueryTable(new(Attachment)).
		OrderBy("-updated_at").
		All(&items); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[TITLE] = Tr(p.Locale(), "nut.attachments.index.title")

	p.TplName = "nut/attachments/index.html"
}

// CreateAttachment upload file
// @router /attachments [post]
func (p *Plugin) CreateAttachment() {
	p.MustSignIn()
	var item FriendLink
	p.Data["item"] = item
	p.Data[TITLE] = Tr(p.Locale(), "buttons.new")
	p.SetSortOrders()
	p.Data["action"] = p.URLFor("nut.Plugin.CreateAdminFriendLink")
	p.TplName = "nut/admin/friend-links/form.html"
}
