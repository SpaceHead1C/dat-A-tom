package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00031, down00031)
}

func up00031(tx *sql.Tx) error {
	query := `-- Registration of send data event function
CREATE FUNCTION set_sent_data(uuid, uuid, char(64), timestamp) RETURNS SETOF json AS $set_sent_data$
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
$set_sent_data$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00031(tx *sql.Tx) error {
	query := `DROP FUNCTION set_sent_data(uuid, uuid, char(64), timestamp);`
	return execQuery(query, tx)
}
