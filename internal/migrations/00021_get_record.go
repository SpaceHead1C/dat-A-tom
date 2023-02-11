package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00021, down00021)
}

func up00021(tx *sql.Tx) error {
	query := `-- Get function for record
CREATE OR REPLACE FUNCTION get_record(uuid) RETURNS SETOF json AS $get_record$
	BEGIN
		RETURN QUERY 
			SELECT
				json_build_object(
					'id', id,
					'reference_type_id', reference_type_id,
					'name', "name",
					'description', description,
					'deletion_mark', deletion_mark,
					'sum', "sum",
					'change_at', change_at::timestamptz
				)
			FROM records
			WHERE id = $1;
	END;
$get_record$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00021(tx *sql.Tx) error {
	query := `DROP FUNCTION IF EXISTS get_record(uuid);`
	return execQuery(query, tx)
}
