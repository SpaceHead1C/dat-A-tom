package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00034, down00034)
}

func up00034(tx *sql.Tx) error {
	query := `-- Changes registration
DO $$ BEGIN
	CREATE FUNCTION value_after_state_change() RETURNS TRIGGER AS $value_after_state_change$
		BEGIN
			INSERT INTO value_changes (record_id, property_id) VALUES (NEW.owner_id, NEW.property_id);
			RETURN NEW;
		END;
	$value_after_state_change$ LANGUAGE plpgsql;

	CREATE TRIGGER t_value_after_state_change AFTER INSERT OR UPDATE ON "values"
		FOR EACH ROW EXECUTE PROCEDURE value_after_state_change();

	CREATE FUNCTION record_after_state_change() RETURNS TRIGGER AS $record_after_state_change$
		BEGIN
			INSERT INTO record_changes (record_id) VALUES (NEW.id);
			RETURN NEW;
		END;
	$record_after_state_change$ LANGUAGE plpgsql;

	CREATE TRIGGER t_record_after_state_change AFTER INSERT OR UPDATE ON records
		FOR EACH ROW EXECUTE PROCEDURE record_after_state_change();

	CREATE FUNCTION property_after_state_change() RETURNS TRIGGER AS $property_after_state_change$
		BEGIN
			INSERT INTO property_changes (property_id) VALUES (NEW.id);
			RETURN NEW;
		END;
	$property_after_state_change$ LANGUAGE plpgsql;

	CREATE TRIGGER t_property_after_state_change AFTER INSERT OR UPDATE ON properties
		FOR EACH ROW EXECUTE PROCEDURE property_after_state_change();

	CREATE FUNCTION reference_type_after_state_change() RETURNS TRIGGER AS $reference_type_after_state_change$
		BEGIN
			INSERT INTO reference_type_changes (reference_type_id) VALUES (NEW.id);
			RETURN NEW;
		END;
	$reference_type_after_state_change$ LANGUAGE plpgsql;

	CREATE TRIGGER t_reference_type_after_state_change AFTER INSERT OR UPDATE ON reference_types
		FOR EACH ROW EXECUTE PROCEDURE reference_type_after_state_change();
END $$;`
	return execQuery(query, tx)
}

func down00034(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	DROP TRIGGER t_reference_type_after_state_change ON reference_types;
	DROP TRIGGER t_property_after_state_change ON properties;
	DROP TRIGGER t_record_after_state_change ON records;
	DROP TRIGGER t_value_after_state_change ON "values";

	DROP FUNCTION reference_type_after_state_change();
	DROP FUNCTION property_after_state_change();
	DROP FUNCTION record_after_state_change();
	DROP FUNCTION value_after_state_change();
END $$;`
	return execQuery(query, tx)
}
