package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00035, down00035)
}

func up00035(tx *sql.Tx) error {
	query := `-- Change types
DO $$ BEGIN
	CREATE TYPE change_types AS ENUM ('ref_type', 'property', 'record', 'value');
EXCEPTION
	WHEN duplicate_object THEN NULL;
END $$;`
	return execQuery(query, tx)
}

func down00035(tx *sql.Tx) error {
	query := `DROP TYPE change_types;`
	return execQuery(query, tx)
}
