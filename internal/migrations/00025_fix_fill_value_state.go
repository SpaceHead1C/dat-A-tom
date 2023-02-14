package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00025, down00025)
}

func up00025(tx *sql.Tx) error {
	query := `
CREATE OR REPLACE FUNCTION value_state_change() RETURNS TRIGGER AS $value_state_change$
	BEGIN
		NEW."sum" = encode(sha256(convert_to(NEW.value::TEXT || '::' || NEW."type" || COALESCE(NEW.reference_type_id::TEXT, ''), 'UTF-8')), 'hex');
		NEW.change_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
$value_state_change$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00025(tx *sql.Tx) error {
	query := `
CREATE OR REPLACE FUNCTION value_state_change() RETURNS TRIGGER AS $value_state_change$
	BEGIN
		NEW."sum" = encode(sha256(convert_to(NEW.value::TEXT || '::' || NEW."type" || COALESCE(NEW.reference_type_id::TEXT, ''), 'UTF-8')), 'hex');
		RETURN NEW;
	END;
$value_state_change$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}
