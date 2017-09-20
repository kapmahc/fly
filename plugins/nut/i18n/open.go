package i18n

import (
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"github.com/kapmahc/fly/plugins/nut/app"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var (
	_items     = make(map[string]string)
	_languages []language.Tag
)

func open() error {
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
		_languages = append(_languages, tag)
		lang := tag.String()
		log.Info("find locale", lang)

		cfg, err := ini.Load(path)
		if err != nil {
			return err
		}

		for _, sec := range cfg.Sections() {
			z := sec.Name()
			for k, v := range sec.KeysHash() {
				log.Debugf("find", z+"."+k)
				_items[z+"."+k] = v
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func init() {
	app.RegisterResource(open)
}
