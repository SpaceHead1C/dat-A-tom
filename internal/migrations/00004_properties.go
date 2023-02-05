package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00004, down00004)
}

func up00004(tx *sql.Tx) error {
	query := `-- Properties
CREATE TABLE IF NOT EXISTS properties (
	id uuid PRIMARY KEY,
	owner_reference_type_id uuid,
	"types" "types"[] NOT NULL,
	reference_type_ids uuid[],
	"name" varchar(128) NOT NULL DEFAULT '',
	description varchar(1024) NOT NULL DEFAULT '',
	CONSTRAINT fk_owner_reference_type
		FOREIGN KEY(owner_reference_type_id)
			REFERENCES reference_types(id)
);`
	return execQuery(query, tx)
}

func down00004(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS properties;`
	return execQuery(query, tx)
}
