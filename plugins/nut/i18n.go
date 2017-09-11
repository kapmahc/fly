package nut

import (
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"golang.org/x/text/language"
)

const (
	// LOCALE locale key
	LOCALE = "locale"
)

//Locale locale
type Locale struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Lang      string    `json:"lang"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName table name
func (*Locale) TableName() string {
	return "locales"
}

// LoadLocales load locales
func LoadLocales() error {
	const ext = ".ini"
	if err := filepath.Walk("locales", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if info.IsDir() || filepath.Ext(name) != ext {
			return err
		}
		tag, err := language.Parse(name[:len(name)-len(ext)])
		if err != nil {
			return err
		}
		lang := tag.String()
		beego.Info("find locale", lang)
		return i18n.SetMessage(lang, path)
	}); err != nil {
		return err
	}
	return nil
}

// Tr translate content to target language.
func Tr(lang, format string, args ...interface{}) string {
	return i18n.Tr(lang, format, args...)
}

func init() {
	beego.AddFuncMap("t", Tr)
}
