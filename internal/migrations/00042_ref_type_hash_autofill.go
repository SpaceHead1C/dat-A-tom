package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00042, down00042)
}

func up00042(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	CREATE OR REPLACE FUNCTION ref_type_state_change() RETURNS TRIGGER AS $ref_type_state_change$
		BEGIN
			NEW."sum" = ref_type_sum(NEW."name", NEW.description);
			RETURN NEW;
		END;
	$ref_type_state_change$ LANGUAGE plpgsql;

	CREATE TRIGGER t_ref_type_state_change BEFORE INSERT OR UPDATE ON reference_types
		FOR EACH ROW EXECUTE PROCEDURE ref_type_state_change();

	ALTER TABLE IF EXISTS reference_types ALTER COLUMN "sum" DROP DEFAULT;

	UPDATE reference_types SET "sum" = ref_type_sum("name", description) WHERE "sum" = '';
END $$;`
	return execQuery(query, tx)
}

func down00042(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	ALTER TABLE reference_types ALTER COLUMN "sum" SET DEFAULT '';

	DROP TRIGGER t_ref_type_state_change ON reference_types;

	DROP FUNCTION ref_type_state_change();
END $$;`
	return execQuery(query, tx)
}
