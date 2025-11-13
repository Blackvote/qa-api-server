package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upSeedData, downSeedData)
}

func upSeedData(tx *sql.Tx) error {
	if _, err := tx.Exec(`
		INSERT INTO questions (text) VALUES
		  ('What is a REST API and what are its core principles?'),
		  ('What is the difference between the HTTP methods PUT and PATCH?'),
		  ('What is a transaction in a database?'),
		  ('What is the purpose of an index in PostgreSQL?');
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO answers (question_id, user_id, text) VALUES
		  (1, '11111111-1111-1111-1111-111111111111', 'A REST API is an architectural style for client-server interaction based on HTTP. Core principles include client-server, statelessness, caching, uniform interface, and separation of concerns.'),
		  (1, '22222222-2222-2222-2222-222222222222', 'REST is not a protocol but a set of constraints. The main idea is to use proper HTTP methods, response codes, and resources instead of “methods” in the URL.'),
		  (2, '11111111-1111-1111-1111-111111111111', 'PUT is usually used for full replacement of a resource, while PATCH is used for partial updates.'),
		  (2, '33333333-3333-3333-3333-333333333333', 'PUT is idempotent by definition, while PATCH is not required to be idempotent.'),
		  (3, '44444444-4444-4444-4444-444444444444', 'A transaction is a logical unit of work that consists of one or more SQL statements and is executed entirely or rolled back as a whole. The guarantees are described by ACID.'),
		  (4, '22222222-2222-2222-2222-222222222222', 'An index speeds up searching for rows by certain columns using an additional data structure. The tradeoff is slower write operations and additional disk space usage.');
	`); err != nil {
		return err
	}

	return nil
}

func downSeedData(tx *sql.Tx) error {
	if _, err := tx.Exec(`DELETE FROM answers WHERE question_id BETWEEN 1 AND 4;`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM questions WHERE id BETWEEN 1 AND 4;`); err != nil {
		return err
	}
	return nil
}
