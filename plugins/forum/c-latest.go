package forum

import (
	"net/http"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// LatestArticles latest
// @router /latest/articles [get]
func (p *Plugin) LatestArticles() {
	p.LayoutApplication()
	var items []Article
	if _, err := orm.NewOrm().QueryTable(new(Article)).
		OrderBy("-updated_at").
		All(&items); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "forum.latest.articles.title")

	p.TplName = "forum/articles/latest.html"
}

// LatestTags latest
// @router /latest/tags [get]
func (p *Plugin) LatestTags() {
	p.LayoutApplication()
	var items []Tag
	if _, err := orm.NewOrm().QueryTable(new(Tag)).
		OrderBy("-updated_at").
		All(&items); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "forum.latest.tags.title")

	p.TplName = "forum/tags/latest.html"
}

// LatestComments latest
// @router /latest/comments [get]
func (p *Plugin) LatestComments() {
	p.LayoutApplication()
	var items []Comment
	if _, err := orm.NewOrm().QueryTable(new(Comment)).
		OrderBy("-updated_at").
		All(&items); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}

	p.Data["items"] = items
	p.Data[nut.TITLE] = nut.Tr(p.Locale(), "forum.latest.comments.title")

	p.TplName = "forum/comments/latest.html"
}
