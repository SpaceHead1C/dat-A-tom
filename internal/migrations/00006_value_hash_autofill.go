package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00006, down00006)
}

func up00006(tx *sql.Tx) error {
	query := `-- Hash autofill
DO $$ BEGIN
	CREATE OR REPLACE FUNCTION value_state_change() RETURNS TRIGGER AS $value_state_change$
		BEGIN
			NEW."sum" = encode(sha256(convert_to(NEW.value::TEXT || '::' || NEW."type" || COALESCE(NEW.reference_type_id::TEXT, ''), 'UTF-8')), 'hex');
			RETURN NEW;
		END;
	$value_state_change$ LANGUAGE plpgsql;

	IF NOT EXISTS (
		SELECT *
		FROM information_schema.triggers
		WHERE event_object_table = 'values'
		AND trigger_name = 't_value_state_change'
	) THEN
		CREATE TRIGGER t_value_state_change BEFORE INSERT OR UPDATE ON "values"
			FOR EACH ROW EXECUTE PROCEDURE value_state_change();
	END IF;
END $$;`
	return execQuery(query, tx)
}

func down00006(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	DROP TRIGGER IF EXISTS t_value_state_change ON "values";
	
	DROP FUNCTION IF EXISTS value_state_change();
END $$;`
	return execQuery(query, tx)
}