package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00032, down00032)
}

func up00032(tx *sql.Tx) error {
	query := `-- Get function for value
CREATE FUNCTION get_value(uuid, uuid) RETURNS SETOF json AS $get_value$
	BEGIN
		RETURN QUERY 
			SELECT
				json_build_object(
					'owner_id', owner_id,
					'property_id', property_id,
					'type', "type",
					'reference_type_id', reference_type_id,
					'value', value,
					'sum', "sum",
					'change_at', change_at::timestamptz
				)
			FROM values
			WHERE owner_id = $1 AND property_id = $2;
	END;
$get_value$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00032(tx *sql.Tx) error {
	query := `DROP FUNCTION get_value(uuid, uuid);`
	return execQuery(query, tx)
}
