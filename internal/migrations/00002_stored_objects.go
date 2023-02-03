package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00002, down00002)
}

func up00002(tx *sql.Tx) error {
	query := `-- Stored objects
CREATE TABLE IF NOT EXISTS records (
	id uuid PRIMARY KEY
	, reference_type_id uuid
	, "name" varchar(128) NOT NULL DEFAULT ''
	, description varchar(1024) NOT NULL DEFAULT ''
	, deletion_mark bool NOT NULL DEFAULT false
	, CONSTRAINT fk_reference_type
			FOREIGN KEY(reference_type_id)
				REFERENCES reference_types(id)
);`
	return execQuery(query, tx)
}

func down00002(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS records;`
	return execQuery(query, tx)
}
