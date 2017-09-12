package nut

import (
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/astaxie/beego"
)

// GetHome home
// @router / [get]
func (p *Plugin) GetHome() {
	p.TplName = "index.tpl"
}

// GetNginxConf nginx.conf
// @router /nginx.conf [get]
func (p *Plugin) GetNginxConf() {
	tpl, err := template.ParseFiles(path.Join("templates", "nginx.conf"))
	if err != nil {
		p.CustomAbort(http.StatusOK, err.Error())
	}
	pwd, _ := os.Getwd()
	ssl, _ := p.GetBool("ssl", false)
	cfg := beego.BConfig
	tpl.Execute(
		p.Ctx.ResponseWriter,
		map[string]interface{}{
			"name":  cfg.ServerName,
			"port":  cfg.Listen.HTTPPort,
			"theme": cfg.AppName,
			"root":  pwd,
			"ssl":   ssl,
		})

}
