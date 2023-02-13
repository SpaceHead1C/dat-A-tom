package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00019, down00019)
}

func up00019(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	CREATE OR REPLACE FUNCTION record_state_change() RETURNS TRIGGER AS $record_state_change$
		BEGIN
			NEW."sum" = record_sum(NEW."name", NEW.description, NEW.deletion_mark);
			NEW.change_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
	$record_state_change$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION property_state_change() RETURNS TRIGGER AS $property_state_change$
		BEGIN
			NEW."sum" = property_sum(NEW."name", NEW.description);
			NEW.change_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
	$property_state_change$ LANGUAGE plpgsql;
END $$;`
	return execQuery(query, tx)
}

func down00019(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	CREATE OR REPLACE FUNCTION property_state_change() RETURNS TRIGGER AS $property_state_change$
		BEGIN
			NEW."sum" = property_sum(NEW."name", NEW.description);
			RETURN NEW;
		END;
	$property_state_change$ LANGUAGE plpgsql;

	CREATE OR REPLACE FUNCTION record_state_change() RETURNS TRIGGER AS $record_state_change$
		BEGIN
			NEW."sum" = record_sum(NEW."name", NEW.description, NEW.deletion_mark);
			RETURN NEW;
		END;
	$record_state_change$ LANGUAGE plpgsql;
END $$;`
	return execQuery(query, tx)
}
