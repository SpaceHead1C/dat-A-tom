package migrations

import (
	"database/sql"
	"datatom/pkg/db/pg"

	"github.com/pressly/goose"
)

func UpMigrations(db *pg.DB) error {
	return goose.Up(db.DB, ".")
}

func execQuery(query string, tx *sql.Tx) error {
	_, err := tx.Exec(query)
	return err
}
