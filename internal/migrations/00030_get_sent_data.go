package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00030, down00030)
}

func up00030(tx *sql.Tx) error {
	query := `-- Get function for sent data
CREATE FUNCTION get_sent_data(uuid, uuid) RETURNS SETOF json AS $get_sent_data$
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
$get_sent_data$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00030(tx *sql.Tx) error {
	query := `DROP FUNCTION get_sent_data(uuid, uuid);`
	return execQuery(query, tx)
}
