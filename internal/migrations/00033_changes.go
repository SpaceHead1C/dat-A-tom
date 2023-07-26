package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00033, down00033)
}

func up00033(tx *sql.Tx) error {
	query := `-- Changes
DO $$ BEGIN
	CREATE SEQUENCE changed_data_id_seq;

	CREATE TABLE reference_type_changes (
		id bigint PRIMARY KEY DEFAULT nextval('changed_data_id_seq'),
		reference_type_id uuid NOT NULL
	);

	CREATE TABLE property_changes (
		id bigint PRIMARY KEY DEFAULT nextval('changed_data_id_seq'),
		property_id uuid NOT NULL
	);

	CREATE TABLE record_changes (
		id bigint PRIMARY KEY DEFAULT nextval('changed_data_id_seq'),
		record_id uuid NOT NULL
	);

	CREATE TABLE value_changes (
		id bigint PRIMARY KEY DEFAULT nextval('changed_data_id_seq'),
		record_id uuid NOT NULL,
		property_id uuid NOT NULL
	);
END $$;`
	return execQuery(query, tx)
}

func down00033(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	DROP TABLE value_changes;
	DROP TABLE record_changes;
	DROP TABLE property_changes;
	DROP TABLE reference_type_changes;

	DROP SEQUENCE changed_data_id_seq;
END $$;`
	return execQuery(query, tx)
}
