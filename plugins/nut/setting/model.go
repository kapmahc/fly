package setting

import "time"

// Model model
type Model struct {
	tableName struct{} `sql:"settings"`
	ID        uint
	Key       string
	Val       []byte
	Encode    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
