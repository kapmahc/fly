package forum

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// Article article
type Article struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	User     *nut.User  `orm:"rel(fk)" json:"user"`
	Tags     []*Tag     `orm:"rel(m2m);rel_table(forum_articles_tags)" json:"tags"`
	Comments []*Comment `orm:"reverse(many)" json:"comments"`
}

// TableName table name
func (*Article) TableName() string {
	return "forum_articles"
}

// Tag tag
type Tag struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Articles []*Article `orm:"reverse(many);rel_table(forum_articles_tags)" json:"articles"`
}

// TableName table name
func (*Tag) TableName() string {
	return "forum_tags"
}

// Comment comment
type Comment struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	User    *nut.User `orm:"rel(fk)" json:"user"`
	Article *Article  `orm:"rel(fk)" json:"article"`
}

// TableName table name
func (*Comment) TableName() string {
	return "forum_comments"
}

func init() {
	orm.RegisterModel(new(Article), new(Tag), new(Comment))
}
