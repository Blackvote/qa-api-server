package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTables, downCreateTables)
}

func upCreateTables(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS questions (
			id SERIAL PRIMARY KEY,
			text TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS answers (
			id SERIAL PRIMARY KEY,
			question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			text TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`); err != nil {
		return err
	}

	return nil
}

func downCreateTables(tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS answers;`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS questions;`); err != nil {
		return err
	}
	return nil
}
