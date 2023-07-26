package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00040, down00040)
}

func up00040(tx *sql.Tx) error {
	query := `-- Sending data event registration functions
DO $$ BEGIN
	ALTER FUNCTION set_sent_data(uuid, uuid, char(64), timestamp) RENAME TO set_sent_value;

	CREATE FUNCTION set_sent_ref_type(uuid, char(64), timestamp) RETURNS SETOF json AS $set_sent_ref_type$
		BEGIN
			RETURN QUERY
				INSERT INTO sent_reference_types (id, "sum", sent_at)
				VALUES ($1, $2, $3)
				ON CONFLICT(id) DO UPDATE SET 
					"sum" = excluded."sum",
					sent_at = excluded.sent_at
				RETURNING json_build_object(
					'id', id,
					'sum', "sum",
					'sent_at', sent_at::timestamptz
				);
		END;
	$set_sent_ref_type$ LANGUAGE plpgsql;

	CREATE FUNCTION set_sent_property(uuid, char(64), timestamp) RETURNS SETOF json AS $set_sent_property$
		BEGIN
			RETURN QUERY
				INSERT INTO sent_properties (id, "sum", sent_at)
				VALUES ($1, $2, $3)
				ON CONFLICT(id) DO UPDATE SET 
					"sum" = excluded."sum",
					sent_at = excluded.sent_at
				RETURNING json_build_object(
					'id', id,
					'sum', "sum",
					'sent_at', sent_at::timestamptz
				);
		END;
	$set_sent_property$ LANGUAGE plpgsql;

	CREATE FUNCTION set_sent_record(uuid, char(64), timestamp) RETURNS SETOF json AS $set_sent_record$
		BEGIN
			RETURN QUERY
				INSERT INTO sent_records (id, "sum", sent_at)
				VALUES ($1, $2, $3)
				ON CONFLICT(id) DO UPDATE SET 
					"sum" = excluded."sum",
					sent_at = excluded.sent_at
				RETURNING json_build_object(
					'id', id,
					'sum', "sum",
					'sent_at', sent_at::timestamptz
				);
		END;
	$set_sent_record$ LANGUAGE plpgsql;
END $$;`
	return execQuery(query, tx)
}

func down00040(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	ALTER FUNCTION set_sent_value(uuid, uuid, char(64), timestamp) RENAME TO set_sent_data;

	DROP FUNCTION set_sent_record(uuid, char(64), timestamp);
	DROP FUNCTION set_sent_property(uuid, char(64), timestamp);
	DROP FUNCTION set_sent_ref_type(uuid, char(64), timestamp);
END $$;`
	return execQuery(query, tx)
}
