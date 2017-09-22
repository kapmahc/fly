package forum

import (
	"time"

	"github.com/kapmahc/fly/plugins/nut"
)

// Article article
type Article struct {
	tableName struct{}  `sql:"forum_articles"`
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	User     *nut.User  `json:"user"`
	Tags     []*Tag     `json:"tags"`
	Comments []*Comment `json:"comments"`
}

// Tag tag
type Tag struct {
	tableName struct{}  `sql:"forum_tags"`
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Articles []*Article `json:"articles"`
}

// Comment comment
type Comment struct {
	tableName struct{}  `sql:"forum_comments"`
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	User    *nut.User `json:"user"`
	Article *Article  `json:"article"`
}
