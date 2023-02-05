package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00012, down00012)
}

func up00012(tx *sql.Tx) error {
	query := `
ALTER TABLE IF EXISTS properties
	ADD COLUMN IF NOT EXISTS "sum" char(64) NOT NULL DEFAULT '',
	ADD COLUMN IF NOT EXISTS change_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;`
	return execQuery(query, tx)
}

func down00012(tx *sql.Tx) error {
	query := `
ALTER TABLE IF EXISTS properties
	DROP COLUMN IF EXISTS "sum",
	DROP COLUMN IF EXISTS change_at;`
	return execQuery(query, tx)
}