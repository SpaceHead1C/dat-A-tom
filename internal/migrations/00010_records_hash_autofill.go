package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00010, down00010)
}

func up00010(tx *sql.Tx) error {
	query := `
ALTER TABLE IF EXISTS records
	ADD COLUMN IF NOT EXISTS "sum" char(64) NOT NULL DEFAULT '',
	ADD COLUMN IF NOT EXISTS change_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;`
	return execQuery(query, tx)
}

func down00010(tx *sql.Tx) error {
	query := `
ALTER TABLE IF EXISTS records
	DROP COLUMN IF EXISTS "sum",
	DROP COLUMN IF EXISTS change_at;`
	return execQuery(query, tx)
}