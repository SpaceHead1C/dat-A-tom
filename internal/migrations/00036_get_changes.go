package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00036, down00036)
}

func up00036(tx *sql.Tx) error {
	query := `-- Get function for changed data
CREATE FUNCTION get_changes() RETURNS TABLE (id bigint, change_type change_types, "key" json) AS $get_changes$
	BEGIN
		RETURN QUERY
			SELECT rtc.id AS id, 'ref_type'::change_types AS change_type, json_build_object('id', rtc.reference_type_id) AS "key"
			FROM reference_type_changes rtc
			UNION ALL
			SELECT pc.id, 'property'::change_types, json_build_object('id', pc.property_id)
			FROM property_changes pc
			UNION ALL
			SELECT rc.id, 'record'::change_types, json_build_object('id', rc.record_id)
			FROM record_changes rc
			UNION ALL
			SELECT vc.id, 'value'::change_types, json_build_object('owner_id', vc.record_id, 'property_id', vc.property_id)
			FROM value_changes vc
			ORDER BY id
			LIMIT 5000;
	END;
$get_changes$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00036(tx *sql.Tx) error {
	query := `DROP FUNCTION get_changes();`
	return execQuery(query, tx)
}
