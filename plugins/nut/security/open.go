package security

import (
	"crypto/aes"

	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/spf13/viper"
)

func open() error {
	key := viper.GetString("secret")
	cip, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	hmacKey = []byte(key)
	aecCip = cip
	return nil
}

func init() {
	app.RegisterResource(open)
}
