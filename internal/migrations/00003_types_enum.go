package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00003, down00003)
}

func up00003(tx *sql.Tx) error {
	query := `-- Data types
DO $$ BEGIN
	CREATE TYPE "types" AS ENUM ('undefined', 'number', 'text', 'bool', 'date', 'uuid', 'ref');
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;`
	return execQuery(query, tx)
}

func down00003(tx *sql.Tx) error {
	query := `DROP TYPE IF EXISTS "types";`
	return execQuery(query, tx)
}
