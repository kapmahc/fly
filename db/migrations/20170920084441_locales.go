package migrations

import (
	"log"

	mig "github.com/go-pg/migrations"
)

func init() {
	const version = "20170920084441_locales"
	mig.Register(func(db mig.DB) error {
		log.Println("migrate database", version)
		_, err := db.Exec(`
			CREATE TABLE locales (
			  id BIGSERIAL PRIMARY KEY,
			  code VARCHAR(255) NOT NULL,
			  lang VARCHAR(8) NOT NULL DEFAULT 'en-US',
			  message TEXT NOT NULL,
			  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
			  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
			);
			CREATE UNIQUE INDEX idx_locales_code_lang ON locales (code, lang);
			CREATE INDEX idx_locales_code ON locales (code);
			CREATE INDEX idx_locales_lang ON locales (lang);
			`)
		return err
	}, func(db mig.DB) error {
		log.Println("rollback database", version)
		_, err := db.Exec(`DROP TABLE locales;`)
		return err
	})
}
