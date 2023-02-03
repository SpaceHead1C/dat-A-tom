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
	return execQuery(query, tx)
}

func down00001(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS reference_types;`
	return execQuery(query, tx)
}
