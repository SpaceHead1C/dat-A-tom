package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00009, down00009)
}

func up00009(tx *sql.Tx) error {
	query := `-- Sent changes
CREATE TABLE IF NOT EXISTS sent (
	record_id uuid NOT NULL,
	property_id uuid NOT NULL,
	"sum" char(64) NOT NULL DEFAULT '???',
	sent_at timestamp NOT NULL DEFAULT '1970-01-01 00:00:00.0',
	CONSTRAINT fk_record
		FOREIGN KEY(record_id)
			REFERENCES records(id),
	CONSTRAINT fk_property
		FOREIGN KEY(property_id)
			REFERENCES properties(id),
	UNIQUE (record_id, property_id)
);`
	return execQuery(query, tx)
}

func down00009(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS sent;`
	return execQuery(query, tx)
}
