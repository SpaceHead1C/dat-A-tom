package main

import (
	"context"
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
	if err != nil {
		panic(err)
	}
	l.Debug(db != nil)
	l.Info("dat(A)tom service is up")
}
