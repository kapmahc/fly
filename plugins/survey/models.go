package survey

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// Form form
type Form struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName table name
func (*Form) TableName() string {
	return "survey_forms"
}

// Field field
type Field struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

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
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Form *Form `orm:"rel(fk)"`
}

// TableName table name
func (*Record) TableName() string {
	return "survey_records"
}
func init() {
	orm.RegisterModel(new(Form), new(Field), new(Record))
}
