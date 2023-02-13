package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00022, down00022)
}

func up00022(tx *sql.Tx) error {
	query := `-- Property constructor
CREATE OR REPLACE FUNCTION new_property(text, text, "types"[], uuid[] DEFAULT NULL, uuid DEFAULT NULL) RETURNS uuid AS $new_property$
	DECLARE
		res uuid;
	BEGIN
		res := uuid_generate_v4();

		INSERT INTO properties (id, "name", description, "types", reference_type_ids, owner_reference_type_id)
		VALUES (res, $1, $2, $3, $4, $5);

		RETURN res;
	END;
$new_property$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00022(tx *sql.Tx) error {
	query := `DROP FUNCTION IF EXISTS new_property(text, text, "types"[], uuid[], uuid);`
	return execQuery(query, tx)
}