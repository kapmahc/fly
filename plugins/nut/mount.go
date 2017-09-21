package nut

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/kapmahc/fly/plugins/nut/i18n"
	"github.com/spf13/viper"
)

func mount(rt *gin.Engine) error {
	ung := rt.Group("/users")
	ung.GET("/sign-in", getUsersSignIn)
	return nil
}

func openRouter(rt *gin.Engine) error {
	theme := viper.GetString("server.theme")
	tpl, err := loadTemplates(filepath.Join("themes", theme, "views"))
	if err != nil {
		return err
	}
	rt.SetHTMLTemplate(tpl)

	rt.Use(i18n.DetectLocale)
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
		}).
		ParseFiles(files...)
}

func init() {
	app.RegisterMount(openRouter, mount)
}
