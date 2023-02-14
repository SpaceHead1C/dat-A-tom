package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00026, down00026)
}

func up00026(tx *sql.Tx) error {
	query := `-- Set function for value
CREATE OR REPLACE FUNCTION set_value(uuid, uuid, "types", uuid, json) RETURNS SETOF json AS $set_value$
	BEGIN
		RETURN QUERY
			INSERT INTO "values" (owner_id, property_id, "type", reference_type_id, value)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT(owner_id, property_id) DO UPDATE SET 
				"type" = excluded."type",
				reference_type_id = excluded.reference_type_id,
				value = excluded.value
			RETURNING json_build_object(
				'owner_id', owner_id,
				'property_id', property_id,
				'type', "type",
				'reference_type_id', reference_type_id,
				'value', value,
				'sum', "sum",
				'change_at', change_at::timestamptz
			);
	END;
$set_value$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00026(tx *sql.Tx) error {
	query := `DROP FUNCTION IF EXISTS set_value(uuid, uuid, "types", uuid, json);`
	return execQuery(query, tx)
}
