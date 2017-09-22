package survey

import (
	"time"

	"github.com/kapmahc/fly/plugins/nut"
)

// Form form
type Form struct {
	tableName struct{}  `sql:"survey_forms"`
	ID        uint      `json:"id"`
	UID       string    `json:"uid"`
	Mode      string    `json:"mode"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	StartUp   time.Time `json:"startUp"`
	ShutDown  time.Time `json:"shutDown"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	User    *nut.User `json:"user"`
	Fields  []*Field  `json:"fields"`
	Records []*Record `json:"records"`
}

// Available available?
func (p *Form) Available() bool {
	now := time.Now()
	return now.After(p.StartUp) && now.Before(p.ShutDown)
}

// Field field
type Field struct {
	tableName struct{}  `sql:"survey_fields"`
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
	Required  bool      `json:"required"`
	SortOrder int       `json:"sortOrder"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Form *Form
}

// Record record
type Record struct {
	tableName struct{}  `sql:"survey_records"`
	ID        uint      `json:"id"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Form *Form
}
