package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00015, down00015)
}

func up00015(tx *sql.Tx) error {
	query := `-- Reference type constructor
CREATE OR REPLACE FUNCTION new_ref_type(text, text) RETURNS uuid AS $new_ref_type$
	DECLARE
		res uuid;
	BEGIN
		res := uuid_generate_v4();

		INSERT INTO reference_types (id, "name", description)
		VALUES (res, $1, $2);

		RETURN res;
	END;
$new_ref_type$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00015(tx *sql.Tx) error {
	query := `DROP IF EXISTS FUNCTION new_ref_type(text, text);`
	return execQuery(query, tx)
}