package nut

// SitemapHandler sitemap handler
type SitemapHandler func() ([]string, error)

var _sitemapHandlers []SitemapHandler

// RegisterSitemap registe sitemap handler
func RegisterSitemap(args ...SitemapHandler) {
	_sitemapHandlers = append(_sitemapHandlers, args...)
}
