package reading

import (
	"time"

	"github.com/kapmahc/fly/plugins/nut"
)

// Book book
type Book struct {
	tableName   struct{}  `sql:"reading_books"`
	ID          uint      `json:"id"`
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
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Note note
type Note struct {
	tableName struct{} `sql:"reading_notes"`
	ID        uint     `json:"id"`
	Type      string
	Body      string
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	User      *nut.User
	Book      *Book
}
