package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00041, down00041)
}

func up00041(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	CREATE FUNCTION ref_type_sum(text, text) RETURNS char(64) AS $ref_type_sum$
		BEGIN
			RETURN encode(sha256(convert_to($1 || '|' || $2, 'UTF-8')), 'hex');
		END;
	$ref_type_sum$ LANGUAGE plpgsql;

	ALTER TABLE IF EXISTS reference_types
		ADD COLUMN "sum" char(64) NOT NULL DEFAULT '',
		ADD COLUMN change_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;

	UPDATE reference_types SET "sum" = ref_type_sum("name", description), change_at = CURRENT_TIMESTAMP;

	DROP FUNCTION get_ref_type(uuid);
	CREATE FUNCTION get_ref_type(uuid) RETURNS SETOF json AS $get_ref_type$
		BEGIN
			RETURN QUERY
				SELECT
					json_build_object(
						'id', id,
						'name', "name",
						'description', description,
						'sum', "sum",
						'change_at', change_at::timestamptz
					)
				FROM reference_types
				WHERE id = $1;
		END;
	$get_ref_type$ LANGUAGE plpgsql;
END $$;`
	return execQuery(query, tx)
}

func down00041(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	DROP FUNCTION get_ref_type(uuid);
	CREATE FUNCTION get_ref_type(uuid) RETURNS SETOF reference_types AS $get_ref_type$
		BEGIN
			RETURN QUERY SELECT id, "name", description
			FROM reference_types
			WHERE id = $1;
		END;
	$get_ref_type$ LANGUAGE plpgsql;

	ALTER TABLE IF EXISTS reference_types
		DROP COLUMN "sum",
		DROP COLUMN change_at;

	DROP FUNCTION ref_type_sum(text, text);
END $$;`
	return execQuery(query, tx)
}
