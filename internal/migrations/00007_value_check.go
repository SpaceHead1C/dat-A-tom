package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00007, down00007)
}

func up00007(tx *sql.Tx) error {
	query := `-- Value check
DO $$ BEGIN
	CREATE OR REPLACE FUNCTION values_row_check_bw() RETURNS TRIGGER AS $values_row_check_bw$
		DECLARE
			pass boolean := FALSE;
		BEGIN
			SELECT INTO pass "types" @> ARRAY[NEW."type"] FROM properties WHERE id = NEW.property_id;
			IF NOT pass THEN
				RAISE EXCEPTION 'unexpected type' USING DETAIL = 'KEYS("values"."type") VALUE(' || NEW."type" || ')';
			END IF;

			IF NEW."type" = 'ref'::types THEN
				IF NEW.reference_type_id IS NULL THEN
					RAISE EXCEPTION 'reference type ID missing' USING DETAIL = 'KEYS("values"."type", "values".reference_type_id) VALUES({' || NEW."type" || ', NULL})';
				END IF;

				SELECT INTO pass reference_type_ids @> ARRAY[NEW.reference_type_id] FROM properties WHERE id = NEW.property_id;
				IF NOT pass THEN
					RAISE EXCEPTION 'unexpected reference type ID' USING DETAIL = 'KEYS("values".reference_type_id) VALUE(' || NEW.reference_type_id || ')';
				END IF;
			ELSIF NEW.reference_type_id IS NOT NULL THEN
				RAISE EXCEPTION 'no need reference type ID cause type is not reference' USING DETAIL = 'KEYS("values"."type", "values".reference_type_id) VALUES({' || NEW."type" || ', ' || NEW.reference_type_id || '})';
			END IF;

			RETURN NEW;
		END;
	$values_row_check_bw$ LANGUAGE plpgsql;
	
	IF NOT EXISTS (
		SELECT *
		FROM information_schema.triggers
		WHERE event_object_table = 'values'
		AND trigger_name = 't_values_row_check_bw'
	) THEN
		CREATE TRIGGER t_values_row_check_bw BEFORE INSERT OR UPDATE ON "values"
			FOR EACH ROW EXECUTE PROCEDURE values_row_check_bw();
	END IF;
END $$;`
	return execQuery(query, tx)
}

func down00007(tx *sql.Tx) error {
	query := `
DO $$ BEGIN
	DROP TRIGGER IF EXISTS t_values_row_check_bw ON "values";
	
	DROP FUNCTION IF EXISTS values_row_check_bw();
END $$;`
	return execQuery(query, tx)
}
