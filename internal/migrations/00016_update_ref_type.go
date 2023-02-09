package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00016, down00016)
}

func up00016(tx *sql.Tx) error {
	query := `-- Update function for reference type
CREATE OR REPLACE FUNCTION update_ref_type(uuid, text, text) RETURNS SETOF reference_types AS $update_ref_type$
	BEGIN
		RETURN QUERY UPDATE reference_types SET
			"name" = COALESCE($2, "name", $2),
			description = COALESCE($3, description, $3)
		WHERE id = $1
		RETURNING *;
	END;
$update_ref_type$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00016(tx *sql.Tx) error {
	query := `DROP IF EXISTS FUNCTION update_ref_type(uuid, text, text);`
	return execQuery(query, tx)
}
