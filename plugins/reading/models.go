package reading

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// Book book
type Book struct {
	ID uint `orm:"column(id)" json:"id"`

	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Lang        string    `json:"lang"`
	File        string    `json:"-"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
	Cover       string    `json:"cover"`
	UpdatedAt   time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt   time.Time `orm:"auto_now_add" json:"createdAt"`
}

// TableName table name
func (*Book) TableName() string {
	return "reading_books"
}

// Note note
type Note struct {
	ID        uint `orm:"column(id)" json:"id"`
	Type      string
	Body      string
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`

	User *nut.User `orm:"rel(fk)"`
	Book *Book     `orm:"rel(fk)"`
}

// TableName table name
func (*Note) TableName() string {
	return "reading_notes"
}
func init() {
	orm.RegisterModel(new(Book), new(Note))
}
