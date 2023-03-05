package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00028, down00028)
}

func up00028(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	-- Set dat(A)way tom id
	CREATE FUNCTION set_config_dataway_tom_id(uuid) RETURNS uuid AS $set_config_dataway_tom_id$
	    DECLARE
			res uuid;
		BEGIN
			UPDATE stored_configs SET dataway_tom_id = $1
			RETURNING dataway_tom_id INTO res;
			 
			RETURN res;
		END;
	$set_config_dataway_tom_id$ LANGUAGE plpgsql;

	-- Get dat(A)way tom id
	CREATE FUNCTION get_config_dataway_tom_id() RETURNS uuid AS $get_config_dataway_tom_id$
		DECLARE
			res uuid;
		BEGIN
			SELECT dataway_tom_id INTO STRICT res
			FROM stored_configs;

			RETURN res;
		END;
	$get_config_dataway_tom_id$ LANGUAGE plpgsql;
END $$;`
	return execQuery(query, tx)
}

func down00028(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	DROP FUNCTION get_config_dataway_tom_id();
	DROP FUNCTION set_config_dataway_tom_id(uuid);
END $$;`
	return execQuery(query, tx)
}
