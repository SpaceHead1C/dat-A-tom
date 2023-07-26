package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00039, down00039)
}

func up00039(tx *sql.Tx) error {
	query := `-- Get functions for sent data
DO $$ BEGIN
	ALTER FUNCTION get_sent_data(uuid, uuid) RENAME TO get_sent_value;

	CREATE FUNCTION get_sent_value_for_update(uuid, uuid) RETURNS SETOF json AS $get_sent_value_for_update$
		BEGIN
			RETURN QUERY 
				SELECT
					json_build_object(
						'record_id', record_id,
						'property_id', property_id,
						'sum', "sum",
						'sent_at', sent_at::timestamptz
					)
				FROM sent_values
				WHERE record_id = $1 AND property_id = $2
				FOR UPDATE OF sent_values;
		END;
	$get_sent_value_for_update$ LANGUAGE plpgsql;

	CREATE FUNCTION get_sent_ref_type_for_update(uuid) RETURNS SETOF json AS $get_sent_ref_type_for_update$
		BEGIN
			RETURN QUERY 
				SELECT
					json_build_object(
						'id', id,
						'sum', "sum",
						'sent_at', sent_at::timestamptz
					)
				FROM sent_reference_types
				WHERE id = $1
				FOR UPDATE OF sent_reference_types;
		END;
	$get_sent_ref_type_for_update$ LANGUAGE plpgsql;

	CREATE FUNCTION get_sent_property_for_update(uuid) RETURNS SETOF json AS $get_sent_property_for_update$
		BEGIN
			RETURN QUERY 
				SELECT
					json_build_object(
						'id', id,
						'sum', "sum",
						'sent_at', sent_at::timestamptz
					)
				FROM sent_properties
				WHERE id = $1
				FOR UPDATE OF sent_properties;
		END;
	$get_sent_property_for_update$ LANGUAGE plpgsql;

	CREATE FUNCTION get_sent_record_for_update(uuid) RETURNS SETOF json AS $get_sent_record_for_update$
		BEGIN
			RETURN QUERY 
				SELECT
					json_build_object(
						'id', id,
						'sum', "sum",
						'sent_at', sent_at::timestamptz
					)
				FROM sent_records
				WHERE id = $1
				FOR UPDATE OF sent_records;
		END;
	$get_sent_record_for_update$ LANGUAGE plpgsql;
END $$;`
	return execQuery(query, tx)
}

func down00039(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	ALTER FUNCTION get_sent_value(uuid, uuid) RENAME TO get_sent_data;

	DROP FUNCTION get_sent_record_for_update(uuid);
	DROP FUNCTION get_sent_property_for_update(uuid);
	DROP FUNCTION get_sent_ref_type_for_update(uuid);
	DROP FUNCTION get_sent_value_for_update(uuid, uuid);
END $$;`
	return execQuery(query, tx)
}
