package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00014, down00014)
}

func up00014(tx *sql.Tx) error {
	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	return execQuery(query, tx)
}

func down00014(tx *sql.Tx) error {
	return nil
}
