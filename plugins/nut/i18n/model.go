package i18n

import (
	"time"

	"github.com/kapmahc/fly/plugins/nut/app"
)

//Model locale model
type Model struct {
	tableName struct{}  `sql:"locales"`
	ID        uint      `json:"id"`
	Lang      string    `json:"lang"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Set set to database
func Set(lang, code, message string) error {
	var it Model
	now := time.Now()
	db := app.DB()
	err := db.Model(&it).
		Column("id", "lang", "code").
		Where("lang = ? AND code = ?", lang, code).
		Select()
	it.Message = message
	it.UpdatedAt = now
	if err == nil {
		_, err := db.Model(&it).Column("message", "updated_att").Update()
		return err
	}
	it.Lang = lang
	it.Code = code
	it.CreatedAt = now
	return db.Insert(&it)
}

// Get get from database
func Get(lang, code string) (string, error) {
	var it Model
	if err := app.DB().Model(&it).
		Column("message").
		Where("lang = ? AND code = ?", lang, code).
		Select(); err != nil {
		return "", err
	}
	return it.Message, nil
}
