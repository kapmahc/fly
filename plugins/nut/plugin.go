package nut

import (
	"errors"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	// TITLE page title's key
	TITLE = "title"
)

// Plugin controller
type Plugin struct {
	Controller
}

func init() {
	beego.AddFuncMap("links", func(loc string) ([]Link, error) {
		var items []Link
		_, err := orm.NewOrm().QueryTable(new(Link)).
			Filter("loc", loc).
			OrderBy("sort_order").
			All(&items, "sort_order", "href", "label")
		return items, err
	})

	beego.AddFuncMap("cards", func(loc string) ([]Card, error) {
		var items []Card
		_, err := orm.NewOrm().QueryTable(new(Card)).
			Filter("loc", loc).
			OrderBy("sort_order").
			All(&items, "title", "href", "logo", "summary", "type", "action")
		return items, err
	})

	beego.AddFuncMap("dtf", func(t time.Time) string {
		return t.Format(time.RFC822)
	})
	beego.AddFuncMap("df", func(t time.Time) string {
		return t.Format(DATE_FORMAT)
	})

	beego.AddFuncMap(
		"dict",
		func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	)

	RegisterSitemap(func() ([]string, error) {
		return []string{
			"/users/sign-in",
			"/users/sign-up",
			"/users/forgot-password",
			"/users/confirm",
			"/users/unlock",
		}, nil
	})
}
