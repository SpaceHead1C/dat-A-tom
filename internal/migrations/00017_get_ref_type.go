package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00017, down00017)
}

func up00017(tx *sql.Tx) error {
	query := `-- Get function for reference type
CREATE OR REPLACE FUNCTION get_ref_type(uuid) RETURNS SETOF reference_types AS $get_ref_type$
	BEGIN
		RETURN QUERY SELECT id, "name", description
		FROM reference_types
		WHERE id = $1;
	END;
$get_ref_type$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00017(tx *sql.Tx) error {
	query := `DROP IF EXISTS FUNCTION get_ref_type(uuid);`
	return execQuery(query, tx)
}
