package i18n

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/kapmahc/fly/plugins/nut/cache"
)

// H html
func H(lang, code string, obj interface{}) (string, error) {
	msg, err := get(lang, code)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("").Parse(msg)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, obj)
	return buf.String(), err
}

//E error
func E(lang, code string, args ...interface{}) error {
	msg, err := get(lang, code)
	if err != nil {
		return errors.New(code)
	}
	return fmt.Errorf(msg, args...)
}

//T text
func T(lang, code string, args ...interface{}) string {
	msg, err := get(lang, code)
	if err != nil {
		return code
	}
	return fmt.Sprintf(msg, args...)
}

func get(lang, code string) (string, error) {
	key := "locales/" + lang + "/" + code
	var msg string
	err := cache.Get(key, &msg)
	if err == nil {
		return msg, nil
	}
	msg, err = Get(lang, code)
	if err == nil {
		cache.Set(key, msg, time.Hour*24)
		return msg, nil
	}

	if msg, ok := _items[lang+"."+code]; ok {
		cache.Set(key, msg, time.Hour*24)
		return msg, nil
	}
	return "", errors.New("not found")
}
