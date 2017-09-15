package nut

import (
	"errors"
	"time"

	"github.com/astaxie/beego"
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
	beego.AddFuncMap("dtf", func(t time.Time) string {
		return t.Format(time.RFC822)
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
}
