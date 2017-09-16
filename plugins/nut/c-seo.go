package nut

import (
	"fmt"
	"net/http"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

// GetGoogleVerify google verify file
// @router /google:code([\w]+).html [get]
func (p *Plugin) GetGoogleVerify() {
	var code string
	if err := Get("site.google.verify-code", &code); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if code != p.Ctx.Input.Param(":code") {
		p.Abort(http.StatusNotFound, nil)
	}
	p.Ctx.WriteString(fmt.Sprintf("google-site-verification: google%s.html", code))
}

// GetBaiduVerify baidu verify file
// @router /baidu_verify_:code([\w]+).html”* [get]
func (p *Plugin) GetBaiduVerify() {
	var code string
	if err := Get("site.baidu.verify-code", &code); err != nil {
		p.Abort(http.StatusInternalServerError, err)
	}
	if code != p.Ctx.Input.Param(":code") {
		p.Abort(http.StatusNotFound, nil)
	}
	p.Ctx.WriteString(code)
}

// GetSitemap sitemap
// @router /sitemap.xml [get]
func (p *Plugin) GetSitemap() {
	sm := stm.NewSitemap()
	sm.SetDefaultHost(p.HomeURL())
	sm.SetCompress(true)
	sm.SetSitemapsPath("/")
	sm.Create()

	for _, l := range p.Data["languages"].([]string) {
		sm.Add(stm.URL{"loc": "/?locale=" + l})
	}

	for _, h := range _sitemapHandlers {
		urls, err := h()
		if err != nil {
			p.Abort(http.StatusInternalServerError, err)
		}
		for _, u := range urls {
			sm.Add(stm.URL{"loc": u})
		}
	}

	p.Ctx.ResponseWriter.Write(sm.XMLContent())
}
