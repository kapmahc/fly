package sitemap

import "github.com/ikeikeikeike/go-sitemap-generator/stm"

// Handler return sitemap's url
type Handler func() ([]stm.URL, error)

var handlers []Handler

// Register register handler
func Register(args ...Handler) {
	handlers = append(handlers, args...)
}
