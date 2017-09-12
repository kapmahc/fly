package nut

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
	"golang.org/x/text/language"
)

//Locale locale
type Locale struct {
	ID        uint      `orm:"column(id)" json:"id"`
	Lang      string    `json:"lang"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	CreatedAt time.Time `orm:"auto_now_add" json:"createdAt"`
	UpdatedAt time.Time `orm:"auto_now" json:"updatedAt"`
}

// TableName table name
func (*Locale) TableName() string {
	return "locales"
}

// SetLocale set locale info in database
func SetLocale(o orm.Ormer, lang, code, message string) error {
	var it Locale
	err := o.QueryTable(&it).
		Filter("lang", lang).
		Filter("code", code).
		One(&it, "id")

	if err == nil {
		_, err = o.QueryTable(&it).Filter("id", it.ID).Update(orm.Params{
			"message":    message,
			"updated_at": time.Now(),
		})
	} else if err == orm.ErrNoRows {
		it.Lang = lang
		it.Code = code
		it.Message = message
		_, err = o.Insert(&it)
	}
	return err
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

// GetLocale get locale message
func GetLocale(lang, code string) (string, error) {
	var it Locale
	if err := orm.NewOrm().QueryTable(&it).
		Filter("lang", lang).
		Filter("code", code).
		One(&it, "Message"); err != nil {
		return "", err
	}
	return it.Message, nil
}

// Th translate content to target language.(html)
func Th(lang, code string, obj interface{}) (string, error) {
	msg, err := GetLocale(lang, code)
	if err != nil {
		msg = i18n.Tr(lang, code)
	}

	tpl, err := template.New("").Parse(msg)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = tpl.Execute(&buf, obj); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Tr translate content to target language.
func Tr(lang, code string, args ...interface{}) string {
	if msg, err := GetLocale(lang, code); err == nil {
		return fmt.Sprintf(msg, args...)
	}
	return i18n.Tr(lang, code, args...)
}

func init() {
	beego.AddFuncMap("t", Tr)
	orm.RegisterModel(new(Locale))
}
