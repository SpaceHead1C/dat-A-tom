package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00029, down00029)
}

func up00029(tx *sql.Tx) error {
	query := `-- Get function for changed values
CREATE FUNCTION get_changed_values() RETURNS SETOF json AS $get_changed_values$
	BEGIN
		RETURN QUERY
			SELECT json_build_object(
				'owner_id', v.owner_id,
				'property_id', v.property_id,
				'type', v."type",
				'reference_type_id', v.reference_type_id,
				'value', v.value,
				'sum', v."sum",
				'change_at', v.change_at::timestamptz
			)
			FROM "values" v LEFT JOIN sent s ON v.owner_id = s.record_id AND v.property_id = s.property_id
			WHERE
				s.sum IS NULL
				OR (v.change_at >= s.sent_at AND v."sum" <> s."sum")
			ORDER BY v.change_at
			LIMIT 5000;
	END;
$get_changed_values$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00029(tx *sql.Tx) error {
	query := `DROP FUNCTION get_changed_values();`
	return execQuery(query, tx)
}
