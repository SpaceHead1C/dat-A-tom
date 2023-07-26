package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up00043, down00043)
}

func up00043(tx *sql.Tx) error {
	query := `-- Purge registration of changed data
CREATE FUNCTION purge_changes(bigint) RETURNS bigint AS $purge_changes$
DECLARE
    res bigint;
BEGIN
    WITH pc AS (
        DELETE FROM property_changes WHERE id <= $1 RETURNING id
    ), rc AS (
        DELETE FROM record_changes WHERE id <= $1 RETURNING id
    ), rtc AS (
        DELETE FROM reference_type_changes WHERE id <= $1 RETURNING id
    ), vc AS (
        DELETE FROM value_changes WHERE id <= $1 RETURNING id
    )
    SELECT sum(r.deleted) INTO STRICT res
    FROM (SELECT count(pc.id) AS deleted FROM pc
         UNION ALL
         SELECT count(rc.id) FROM rc
         UNION ALL
         SELECT count(rtc.id) FROM rtc
         UNION ALL
         SELECT count(vc.id) FROM vc) r;
    RETURN res;
END;
$purge_changes$ LANGUAGE plpgsql;`
	return execQuery(query, tx)
}

func down00043(tx *sql.Tx) error {
	query := `DROP FUNCTION purge_changes(bigint);`
	return execQuery(query, tx)
}
