package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00023, down00023)
}

func up00023(tx *sql.Tx) error {
	query := `-- Update function for property
CREATE OR REPLACE FUNCTION update_property(uuid, text, text) RETURNS SETOF json AS $update_property$
	BEGIN
		RETURN QUERY UPDATE properties SET
			"name" = COALESCE($2, "name", $2),
			description = COALESCE($3, description, $3)
		WHERE id = $1
		RETURNING json_build_object(
			'id', id,
			'name', "name",
			'description', description,
			'types', "types",
			'reference_type_ids', reference_type_ids,
			'owner_reference_type_id', owner_reference_type_id,
			'sum', "sum",
			'change_at', change_at::timestamptz
		);
	END;
$update_property$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00023(tx *sql.Tx) error {
	query := `DROP FUNCTION IF EXISTS update_property(uuid, text, text);`
	return execQuery(query, tx)
}
