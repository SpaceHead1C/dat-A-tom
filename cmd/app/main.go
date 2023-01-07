package main

import (
	"context"
	"datatom/pkg/db/pg"
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
	db, err := pg.NewDB(dbCtx, pg.Config{
		Address:      c.PostgresAddress,
		Port:         c.PostgresPort,
		User:         c.PostgresUser,
		Password:     c.PostgresPassword,
		DatabaseName: c.PostgresDBName,
	})
	cancel()
	if err != nil {
		panic(err)
	}
	l.Debug(db != nil)
	l.Info("dat(A)tom service is up")
}
