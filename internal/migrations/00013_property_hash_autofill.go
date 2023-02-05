package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00013, down00013)
}

func up00013(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	CREATE OR REPLACE FUNCTION property_sum(text, text) RETURNS char(64) AS $property_sum$
		BEGIN
			RETURN encode(sha256(convert_to($1 || '|' || $2, 'UTF-8')), 'hex');
		END;
	$property_sum$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION property_state_change() RETURNS TRIGGER AS $property_state_change$
		BEGIN
			NEW."sum" = property_sum(NEW."name", NEW.description);
			RETURN NEW;
		END;
	$property_state_change$ LANGUAGE plpgsql;

	IF NOT EXISTS (
		SELECT *
		FROM information_schema.triggers
		WHERE event_object_table = 'properties'
		AND trigger_name = 't_property_state_change'
	) THEN
		CREATE TRIGGER t_property_state_change BEFORE INSERT OR UPDATE ON properties
			FOR EACH ROW EXECUTE PROCEDURE property_state_change();
	END IF;

	ALTER TABLE IF EXISTS properties ALTER COLUMN "sum" DROP DEFAULT;

	UPDATE properties SET "sum" = property_sum("name", description) WHERE "sum" = '';
END $$;`
	return execQuery(query, tx)
}

func down00013(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	ALTER TABLE IF EXISTS properties ALTER COLUMN "sum" SET DEFAULT '';

	DROP TRIGGER IF EXISTS t_property_state_change ON records;

	DROP FUNCTION IF EXISTS property_state_change();
	DROP FUNCTION IF EXISTS property_sum();
END $$;`
	return execQuery(query, tx)
}