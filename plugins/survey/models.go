package survey

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/kapmahc/fly/plugins/nut"
)

// Form form
type Form struct {
	ID        uint      `orm:"column(id)" json:"id"`
	UID       string    `orm:"column(uid)" json:"uid"`
	Mode      string    `json:"mode"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	StartUp   time.Time `json:"startUp"`
	ShutDown  time.Time `json:"shutDown"`
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`

	User    *nut.User `orm:"rel(fk)" json:"user"`
	Fields  []*Field  `orm:"reverse(many)" json:"fields"`
	Records []*Record `orm:"reverse(many)" json:"records"`
}

// TableName table name
func (*Form) TableName() string {
	return "survey_forms"
}

// Available available?
func (p *Form) Available() bool {
	now := time.Now()
	return now.After(p.StartUp) && now.Before(p.ShutDown)
}

// Field field
type Field struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
	Required  bool      `json:"required"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`

	Form *Form `orm:"rel(fk)"`
}

// TableName table name
func (*Field) TableName() string {
	return "survey_fields"
}

// Record record
type Record struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`

	Form *Form `orm:"rel(fk)"`
}

// TableName table name
func (*Record) TableName() string {
	return "survey_records"
}
func init() {
	orm.RegisterModel(new(Form), new(Field), new(Record))
}
