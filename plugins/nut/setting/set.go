package setting

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/kapmahc/fly/plugins/nut/security"
)

// Set set k-v
func Set(k string, v interface{}, f bool) error {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return err
	}
	val := buf.Bytes()
	if f {
		var err error
		val, err = security.Encrypt(val)
		if err != nil {
			return err
		}
	}

	var it Model
	now := time.Now()
	db := app.DB()
	err := db.Model(&it).
		Column("id", "key", "val").
		Where("key = ?", k).
		Select()
	it.Encode = f
	it.Val = val
	it.UpdatedAt = now
	if err == nil {
		_, err := db.Model(&it).Column("encode", "val", "updated_att").Update()
		return err
	}
	it.Key = k
	it.CreatedAt = now
	return db.Insert(&it)
}
