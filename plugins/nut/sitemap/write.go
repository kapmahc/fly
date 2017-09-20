package sitemap

import (
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/kapmahc/fly/plugins/nut/app"
)

// Write write sitemap.xml.gz to tmp
func Write(ping bool) error {
	sm := stm.NewSitemap()
	sm.SetDefaultHost(app.Home())
	sm.SetPublicPath("tmp/")
	sm.SetSitemapsPath("/")
	sm.SetCompress(true)
	sm.SetVerbose(true)
	sm.Create()

	for _, hnd := range handlers {
		urls, err := hnd()
		if err != nil {
			return err
		}
		for _, u := range urls {
			sm.Add(u)
		}
	}
	sm.Finalize()
	return nil
}
