package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00038, down00038)
}

func up00038(tx *sql.Tx) error {
	query := `-- Sent changes
DO $$ BEGIN
	CREATE TABLE sent_reference_types (
		id uuid NOT NULL,
		"sum" char(64) NOT NULL DEFAULT '???',
		sent_at timestamp NOT NULL DEFAULT '1970-01-01 00:00:00.0',
		CONSTRAINT fk_reference_type
			FOREIGN KEY (id)
				REFERENCES reference_types(id),
		UNIQUE (id)
	);
	
	CREATE TABLE sent_properties (
		id uuid NOT NULL,
		"sum" char(64) NOT NULL DEFAULT '???',
		sent_at timestamp NOT NULL DEFAULT '1970-01-01 00:00:00.0',
		CONSTRAINT fk_property
			FOREIGN KEY (id)
				REFERENCES properties(id),
		UNIQUE (id)
	);

	CREATE TABLE sent_records (
		id uuid NOT NULL,
		"sum" char(64) NOT NULL DEFAULT '???',
		sent_at timestamp NOT NULL DEFAULT '1970-01-01 00:00:00.0',
		CONSTRAINT fk_record
			FOREIGN KEY (id)
				REFERENCES records(id),
		UNIQUE (id)
	);
END $$;`
	return execQuery(query, tx)
}

func down00038(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	DROP TABLE sent_records;
	DROP TABLE sent_properties;
	DROP TABLE sent_reference_types;
END $$;`
	return execQuery(query, tx)
}
