package rss

import "golang.org/x/tools/blog/atom"

// Handler return sitemap's url
type Handler func(string) ([]*atom.Entry, error)

var handlers []Handler

// Register register handler
func Register(args ...Handler) {
	handlers = append(handlers, args...)
}
