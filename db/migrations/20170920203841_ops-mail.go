package migrations

import (
	"log"

	mig "github.com/go-pg/migrations"
)

func init() {
	const version = "20170920203841_ops-mail"
	mig.Register(func(db mig.DB) error {
		log.Println("migrate database", version)
		_, err := db.Exec(`
			CREATE TABLE mail_domains (
			  id         BIGSERIAL PRIMARY KEY,
			  name       VARCHAR(128)                NOT NULL,
			  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
			  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
			);
			CREATE UNIQUE INDEX idx_mail_domains_name
			  ON mail_domains (name);

			CREATE TABLE mail_users (
			  id         BIGSERIAL PRIMARY KEY,
			  domain_id  BIGINT                      NOT NULL REFERENCES mail_domains,
			  email      VARCHAR(255)                NOT NULL,
			  full_name  VARCHAR(128)                NOT NULL,
			  password   VARCHAR(255)                NOT NULL,
			  enable     BOOLEAN  NOT NULL,
			  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
			  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
			);
			CREATE UNIQUE INDEX idx_mail_users_email
			  ON mail_users (email);
			CREATE INDEX idx_mail_users_full_name
			  ON mail_users (full_name);

			CREATE TABLE mail_aliases (
			  id          BIGSERIAL PRIMARY KEY,
			  domain_id   BIGINT                      NOT NULL REFERENCES mail_domains,
			  source      VARCHAR(255)                NOT NULL,
			  destination VARCHAR(255)                NOT NULL REFERENCES mail_users(email),
			  created_at  TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
			  updated_at  TIMESTAMP WITHOUT TIME ZONE NOT NULL
			);
			CREATE UNIQUE INDEX idx_mail_aliases_source
			  ON mail_aliases (source);
			CREATE INDEX idx_mail_aliases_destination
			  ON mail_aliases (destination);				
			`)
		return err
	}, func(db mig.DB) error {
		log.Println("rollback database", version)
		_, err := db.Exec(`
			DROP TABLE mail_aliases;
			DROP TABLE mail_users;
			DROP TABLE mail_domains;
			`)
		return err
	})
}