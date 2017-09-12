package nut

import (
	"crypto/aes"

	"github.com/astaxie/beego"
)

// Open open
func Open() error {
	if err := LoadLocales(); err != nil {
		return err
	}
	// init security
	cip, err := aes.NewCipher([]byte(beego.AppConfig.String("aeskey")))
	if err != nil {
		return err
	}
	_aes = &Aes{cip: cip}
	_hmac = &Hmac{key: []byte(beego.AppConfig.String("hmackey"))}
	return nil
}
