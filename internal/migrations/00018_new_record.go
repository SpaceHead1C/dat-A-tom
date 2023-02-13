package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00018, down00018)
}

func up00018(tx *sql.Tx) error {
	query := `-- Record constructor
CREATE OR REPLACE FUNCTION new_record(text, text, bool DEFAULT FALSE, uuid DEFAULT NULL) RETURNS uuid AS $new_record$
	DECLARE
		res uuid;
	BEGIN
		res := uuid_generate_v4();

		INSERT INTO records (id, "name", description, deletion_mark, reference_type_id)
		VALUES (res, $1, $2, $3, $4);

		RETURN res;
	END;
$new_record$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00018(tx *sql.Tx) error {
	query := `DROP FUNCTION IF EXISTS new_record(text, text, bool, uuid);`
	return execQuery(query, tx)
}
