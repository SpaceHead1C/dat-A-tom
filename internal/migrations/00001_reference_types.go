package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00001, down00001)
}

func up00001(tx *sql.Tx) error {
	query := `-- Reference types
CREATE TABLE IF NOT EXISTS reference_types (
	id uuid PRIMARY KEY,
	"name" varchar(128) NOT NULL DEFAULT '',
	description varchar(1024) NOT NULL DEFAULT ''
);`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}

func down00001(tx *sql.Tx) error {
	query := `DROP TABLE reference_types;`
	if _, err := tx.Exec(query); err != nil {
		return err
	}
	return nil
}
