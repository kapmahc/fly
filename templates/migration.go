package migrations

import (
	"log"

	mig "github.com/go-pg/migrations"
)

func init() {
	const version = "{{.Version}}_{{.Name}}"
	mig.Register(func(db mig.DB) error {
		log.Println("migrate database", version)
		_, err := db.Exec(`TODO`)
		return err
	}, func(db mig.DB) error {
		log.Println("rollback database", version)
		_, err := db.Exec(`TODO`)
		return err
	})
}
