package main

import (
	"context"
	"datatom/internal/adapter/pg"
	"datatom/internal/migrations"
	pkgpg "datatom/pkg/db/pg"
	"datatom/pkg/log"
	"os"
	"time"
)

func main() {
	c := newConfig()
	if err := parse(os.Args[1:], c); err != nil {
		panic(err.Error())
	}
	l, err := log.NewLogger()
	if err != nil {
		panic(err.Error())
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	db, err := pkgpg.NewDB(dbCtx, pkgpg.Config{
		Address:      c.PostgresAddress,
		Port:         c.PostgresPort,
		User:         c.PostgresUser,
		Password:     c.PostgresPassword,
		DatabaseName: c.PostgresDBName,
	})
	if err != nil {
		cancel()
		panic(err.Error())
	}
	cancel()
	defer db.Close()
	l.Info("database connected")

	if err := migrations.UpMigrations(db); err != nil {
		panic(err.Error())
	}
	l.Info("migrations got up")

	repo, err := pg.NewRepository(db, l)
	if err != nil {
		panic(err.Error())
	}
	defer repo.CloseConn(db)
	l.Info("repository configured")

	l.Info("dat(A)tom service is up")
}
