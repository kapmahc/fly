package main

import (
	"log"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		log.Println("migrate {{.Version}}_{{.Name}}")
		_, err := db.Exec(`TODO`)
		return err
	}, func(db migrations.DB) error {
		log.Println("rollback {{.Version}}_{{.Name}}")
		_, err := db.Exec(`TODO`)
		return err
	})
}
