package setting

import (
	"bytes"
	"encoding/gob"

	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/kapmahc/fly/plugins/nut/security"
)

// Get get val by key
func Get(k string, v interface{}) error {
	var it Model
	if err := app.DB().Model(&it).
		Column("val").
		Where("key = ?", k).
		Select(); err != nil {
		return err
	}

	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	if it.Encode {
		val, err := security.Decrypt(it.Val)
		if err != nil {
			return err
		}
		buf.Write(val)
	} else {
		buf.Write(it.Val)
	}
	return dec.Decode(v)
}
