package migrations

import (
	"log"

	mig "github.com/go-pg/migrations"
)

func init() {
	const version = "20170920084516_settings"
	mig.Register(func(db mig.DB) error {
		log.Println("migrate database", version)
		_, err := db.Exec(`
			CREATE TABLE settings (
			  id BIGSERIAL PRIMARY KEY,
			  key VARCHAR(255) NOT NULL,
			  val BYTEA NOT NULL,
			  encode BOOLEAN NOT NULL DEFAULT FALSE,
			  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
			  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
			);
			CREATE UNIQUE INDEX idx_settings_key ON settings (key);
			`)
		return err
	}, func(db mig.DB) error {
		log.Println("rollback database", version)
		_, err := db.Exec(`DROP TABLE settings;`)
		return err
	})
}
