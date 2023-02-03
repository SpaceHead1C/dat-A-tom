package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00008, down00008)
}

func up00008(tx *sql.Tx) error {
	query := `-- Property check
DO $$ BEGIN
	CREATE OR REPLACE FUNCTION properties_row_check_bw() RETURNS TRIGGER AS $properties_row_check_bw$
		DECLARE
			fail boolean := FALSE;
			vtps TEXT := '';
			vids TEXT := '';
		BEGIN
			IF CARDINALITY(NEW."types") = 0 THEN
				IF NEW."types" IS NULL THEN
					vtps := 'NULL';
				ELSE
					vtps := '{' || array_to_string(NEW."types", ', ') || '}';
				END IF;
				RAISE EXCEPTION 'types expected' USING DETAIL = 'KEYS(properties."types") VALUES(' || vtps || ')';
			END IF;
			
			IF (array_position(NEW."types", 'ref') IS NULL)::int # (CARDINALITY(COALESCE(NEW.reference_type_ids, '{}'::uuid[])) = 0)::int > 0 THEN
				IF NEW."types" IS NULL THEN
					vtps := 'NULL';
				ELSE
					vtps := '{' || array_to_string(NEW."types", ', ') || '}';
				END IF;
				IF NEW.reference_type_ids IS NULL THEN
					vids := 'NULL';
				ELSE
					vids := '{' || array_to_string(NEW.reference_type_ids, ', ') || '}';
				END IF;
				RAISE EXCEPTION 'types and reference type condition not matched' USING DETAIL = 'KEYS(properties."types", properties.reference_type_ids) VALUES(' || vtps || ', ' || vids || ')';
			END IF;
			
			SELECT EXISTS INTO fail (SELECT u FROM UNNEST(NEW."types") u GROUP BY u HAVING count(u) > 1);
			IF fail THEN
				RAISE EXCEPTION 'type duplicated' USING DETAIL = 'KEYS(properties."types") VALUES({' || array_to_string(NEW."types", ', ') || '})';
			END IF;
			
			SELECT EXISTS INTO fail (SELECT u FROM UNNEST(NEW.reference_type_ids) u GROUP BY u HAVING count(u) > 1);
			IF fail THEN
				RAISE EXCEPTION 'reference type ID duplicated' USING DETAIL = 'KEYS(properties.reference_type_ids) VALUES({' || array_to_string(NEW.reference_type_ids, ', ') || '})';
			END IF;
			
			SELECT EXISTS INTO fail (SELECT * FROM (SELECT u FROM UNNEST(NEW.reference_type_ids) u) r LEFT JOIN reference_types rt ON r.u = rt.id WHERE rt.id IS NULL);
			IF fail THEN
				RAISE EXCEPTION 'unknown reference type ID' USING DETAIL = 'KEYS(properties.reference_type_ids) VALUES({' || array_to_string(NEW.reference_type_ids, ', ') || '})';
			END IF;
			
			RETURN NEW;
		END;
	$properties_row_check_bw$ LANGUAGE plpgsql;

	IF NOT EXISTS (
		SELECT *
		FROM information_schema.triggers
		WHERE event_object_table = 'properties'
		AND trigger_name = 't_properties_row_check_bw'
	) THEN
		CREATE TRIGGER t_properties_row_check_bw BEFORE INSERT OR UPDATE ON properties
			FOR EACH ROW EXECUTE PROCEDURE properties_row_check_bw();
	END IF;
END $$;`
	return execQuery(query, tx)
}

func down00008(tx *sql.Tx) error {
	query := `
	DO $$ BEGIN
		DROP TRIGGER IF EXISTS t_properties_row_check_bw ON properties;
		
		DROP FUNCTION IF EXISTS properties_row_check_bw();
	END $$;`
	return execQuery(query, tx)
}
