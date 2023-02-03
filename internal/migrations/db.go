package migrations

import (
	"datatom/pkg/db/pg"

	"github.com/pressly/goose"
)

func UpMigrations(db *pg.DB) error {
	return goose.Up(db.DB, ".")
}
