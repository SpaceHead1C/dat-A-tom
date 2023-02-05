package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00011, down00011)
}

func up00011(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	CREATE OR REPLACE FUNCTION record_sum(text, text, bool) RETURNS char(64) AS $record_sum$
		BEGIN
			RETURN encode(sha256(convert_to($1 || '|' || $2 || '|' || $3, 'UTF-8')), 'hex');
		END;
	$record_sum$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION record_state_change() RETURNS TRIGGER AS $record_state_change$
		BEGIN
			NEW."sum" = record_sum(NEW."name", NEW.description, NEW.deletion_mark);
			RETURN NEW;
		END;
	$record_state_change$ LANGUAGE plpgsql;

	IF NOT EXISTS (
		SELECT *
		FROM information_schema.triggers
		WHERE event_object_table = 'records'
		AND trigger_name = 't_record_state_change'
	) THEN
		CREATE TRIGGER t_record_state_change BEFORE INSERT OR UPDATE ON records
			FOR EACH ROW EXECUTE PROCEDURE record_state_change();
	END IF;

	ALTER TABLE IF EXISTS records ALTER COLUMN "sum" DROP DEFAULT;

	UPDATE records SET "sum" = record_sum("name", description, deletion_mark) WHERE "sum" = '';
END $$;`
	return execQuery(query, tx)
}

func down00011(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	ALTER TABLE IF EXISTS records ALTER COLUMN "sum" SET DEFAULT '';

	DROP TRIGGER IF EXISTS t_record_state_change ON records;

	DROP FUNCTION IF EXISTS record_state_change();
	DROP FUNCTION IF EXISTS record_sum();
END $$;`
	return execQuery(query, tx)
}