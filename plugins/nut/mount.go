package nut

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/kapmahc/fly/plugins/nut/i18n"
	"github.com/spf13/viper"
)

func mount(rt *gin.Engine) error {
	ung := rt.Group("/users")
	ung.GET("/sign-in", getUsersSignIn)

	rt.GET("/install", HTML("nut.install", getInstall))
	rt.POST("/install", FORM("/install", &fmInstall{}, postInstall))
	return nil
}

func openRouter(rt *gin.Engine) error {
	theme := viper.GetString("server.theme")
	tpl, err := loadTemplates(filepath.Join("themes", theme, "views"))
	if err != nil {
		return err
	}
	rt.SetHTMLTemplate(tpl)

	rt.Use(
		sessions.Sessions(
			"session",
			sessions.NewCookieStore([]byte(viper.GetString("secret"))),
		),
	)
	rt.Use(i18n.DetectLocale)

	rt.Static("/3rd", "node_modules")
	rt.Static("/assets", filepath.Join("themes", theme, "assets"))
	return nil
}

func loadTemplates(root string) (*template.Template, error) {
	var files []string
	if err := filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			files = append(files, path)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return template.New("").
		Funcs(template.FuncMap{
			"t": i18n.T,
			"assets_js": func(u string) template.HTML {
				return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, u))
			},
			"assets_css": func(u string) template.HTML {
				return template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s">`, u))
			},
			"eq": func(a, b interface{}) bool {
				return a == b
			},
			"printf": fmt.Sprintf,
			"str2htm": func(s string) template.HTML {
				return template.HTML(s)
			},
			"links": func(loc string) []interface{} {
				// TODO
				return []interface{}{}
			},
		}).
		ParseFiles(files...)
}

func init() {
	app.RegisterMount(openRouter, mount)
}
