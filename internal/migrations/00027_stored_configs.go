package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00027, down00027)
}

func up00027(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	-- Stored configs
	CREATE TABLE stored_configs (
		dataway_tom_id uuid
	);
	
	INSERT INTO stored_configs DEFAULT VALUES;
END $$;`
	return execQuery(query, tx)
}

func down00027(tx *sql.Tx) error {
	query := `DROP TABLE stored_configs;`
	return execQuery(query, tx)
}
