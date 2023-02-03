package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00005, down00005)
}

func up00005(tx *sql.Tx) error {
	query := `-- Values
DO $$ BEGIN
	CREATE TABLE IF NOT EXISTS "values" (
		owner_id uuid NOT NULL,
		property_id uuid NOT NULL,
		"type" "types" NOT NULL,
		reference_type_id uuid,
		value jsonb NOT NULL,
		"sum" char(64) NOT NULL,
		change_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(owner_id, property_id),
		CONSTRAINT fk_owner_records
			FOREIGN KEY(owner_id) 
				REFERENCES records(id),
		CONSTRAINT fk_property
			FOREIGN KEY(property_id)
				REFERENCES properties(id),
		CONSTRAINT fk_reference_type
			FOREIGN KEY(reference_type_id)
				REFERENCES reference_types(id)
	);
	
	CREATE INDEX IF NOT EXISTS values_property_idx ON "values" (property_id);

	CREATE INDEX IF NOT EXISTS values_owner_idx ON "values" (owner_id);
END $$;`
	return execQuery(query, tx)
}

func down00005(tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS "values";`
	return execQuery(query, tx)
}
