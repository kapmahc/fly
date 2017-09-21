package migrations

import (
	"log"

	mig "github.com/go-pg/migrations"
)

func init() {
	const version = "20170920203832_reading"
	mig.Register(func(db mig.DB) error {
		log.Println("migrate database", version)
		_, err := db.Exec(`
			CREATE TABLE reading_books (
			  id           BIGSERIAL PRIMARY KEY,
			  author       VARCHAR(255)                NOT NULL,
			  publisher    VARCHAR(255)                NOT NULL,
			  title        VARCHAR(255)                NOT NULL,
			  type         VARCHAR(36)                 NOT NULL DEFAULT 'epub3',
			  lang         VARCHAR(32)                 NOT NULL DEFAULT 'en-US',
			  file         VARCHAR(255)                NOT NULL,
			  subject      VARCHAR(255),
			  description  TEXT,
			  published_at DATE                        NOT NULL DEFAULT current_date,
			  cover        VARCHAR(255),
			  created_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
			  updated_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL
			);
			CREATE UNIQUE INDEX idx_reading_books_file
			  ON reading_books (file);
			CREATE INDEX idx_reading_books_author
			  ON reading_books (author);
			CREATE INDEX idx_reading_books_publisher
			  ON reading_books (publisher);
			CREATE INDEX idx_reading_books_type
			  ON reading_books (type);
			CREATE INDEX idx_reading_books_lang
			  ON reading_books (lang);

			CREATE TABLE reading_notes (
			  id         BIGSERIAL PRIMARY KEY,
			  user_id    BIGINT                      NOT NULL REFERENCES users,
			  book_id    BIGINT                      NOT NULL REFERENCES reading_books,
			  body       TEXT                        NOT NULL,
			  type       VARCHAR(8)                  NOT NULL DEFAULT 'markdown',
			  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
			  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
			);
			CREATE INDEX idx_reading_notes_type
			  ON reading_notes (type);
			`)
		return err
	}, func(db mig.DB) error {
		log.Println("rollback database", version)
		_, err := db.Exec(`
			DROP TABLE reading_notes;
			DROP TABLE reading_books;
			`)
		return err
	})
}
