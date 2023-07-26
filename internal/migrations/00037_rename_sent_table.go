package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00037, down00037)
}

func up00037(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	ALTER TABLE sent RENAME TO sent_values;

	CREATE OR REPLACE FUNCTION get_sent_data(uuid, uuid) RETURNS SETOF json AS $get_sent_data$
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
				WHERE record_id = $1 AND property_id = $2;
		END;
	$get_sent_data$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION set_sent_data(uuid, uuid, char(64), timestamp) RETURNS SETOF json AS $set_sent_data$
		BEGIN
			RETURN QUERY
				INSERT INTO sent_values (record_id, property_id, "sum", sent_at)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT(record_id, property_id) DO UPDATE SET 
					"sum" = excluded."sum",
					sent_at = excluded.sent_at
				RETURNING json_build_object(
					'record_id', record_id,
					'property_id', property_id,
					'sum', "sum",
					'sent_at', sent_at::timestamptz
				);
		END;
	$set_sent_data$ LANGUAGE plpgsql;
END $$;`
	return execQuery(query, tx)
}

func down00037(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	ALTER TABLE sent_values RENAME TO sent;

	CREATE OR REPLACE FUNCTION get_sent_data(uuid, uuid) RETURNS SETOF json AS $get_sent_data$
		BEGIN
			RETURN QUERY 
				SELECT
					json_build_object(
						'record_id', record_id,
						'property_id', property_id,
						'sum', "sum",
						'sent_at', sent_at::timestamptz
					)
				FROM sent
				WHERE record_id = $1 AND property_id = $2;
		END;
	$get_sent_data$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION set_sent_data(uuid, uuid, char(64), timestamp) RETURNS SETOF json AS $set_sent_data$
		BEGIN
			RETURN QUERY
				INSERT INTO sent (record_id, property_id, "sum", sent_at)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT(record_id, property_id) DO UPDATE SET 
					"sum" = excluded."sum",
					sent_at = excluded.sent_at
				RETURNING json_build_object(
					'record_id', record_id,
					'property_id', property_id,
					'sum', "sum",
					'sent_at', sent_at::timestamptz
				);
		END;
	$set_sent_data$ LANGUAGE plpgsql;
END $$;`
	return execQuery(query, tx)
}
