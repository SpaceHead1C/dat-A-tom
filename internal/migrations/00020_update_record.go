package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00020, down00020)
}

func up00020(tx *sql.Tx) error {
	query := `-- Update function for record
CREATE OR REPLACE FUNCTION update_record(uuid, text, text, bool) RETURNS SETOF json AS $update_record$
	BEGIN
		RETURN QUERY UPDATE records SET
			"name" = COALESCE($2, "name", $2),
			description = COALESCE($3, description, $3),
			deletion_mark = COALESCE($4, deletion_mark, $4)
		WHERE id = $1
		RETURNING json_build_object(
			'id', id,
			'reference_type_id', reference_type_id,
			'name', "name",
			'description', description,
			'deletion_mark', deletion_mark,
			'sum', "sum",
			'change_at', change_at::timestamptz
		);
	END;
$update_record$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00020(tx *sql.Tx) error {
	query := `DROP FUNCTION IF EXISTS update_record(uuid, text, text, bool);`
	return execQuery(query, tx)
}
