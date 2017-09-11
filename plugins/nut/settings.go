package nut

import "time"

// Setting k-v
type Setting struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Key       string    `json:"key"`
	Val       []byte    `json:"val"`
	Encode    bool      `json:"encode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (*Setting) TableName() string {
	return "settings"
}
