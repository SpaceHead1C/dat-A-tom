package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00024, down00024)
}

func up00024(tx *sql.Tx) error {
	query := `-- Get function for property
CREATE OR REPLACE FUNCTION get_property(uuid) RETURNS SETOF json AS $get_property$
	BEGIN
		RETURN QUERY
			SELECT
				json_build_object(
					'id', id,
					'name', "name",
					'description', description,
					'types', "types",
					'reference_type_ids', reference_type_ids,
					'owner_reference_type_id', owner_reference_type_id,
					'sum', "sum",
					'change_at', change_at::timestamptz
				)
			FROM properties
			WHERE id = $1;
	END;
$get_property$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00024(tx *sql.Tx) error {
	query := `DROP FUNCTION IF EXISTS get_property(uuid);`
	return execQuery(query, tx)
}
