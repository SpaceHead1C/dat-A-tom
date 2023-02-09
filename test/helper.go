package test

import (
	"context"
	"datatom/internal/api"
	apg "datatom/internal/adapter/pg"
	pkgpg "datatom/pkg/db/pg"
	"datatom/pkg/log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/subosito/gotenv"
)

func newPgRepo(t *testing.T) *apg.Repository {
	if err := gotenv.Load(); err != nil {
		t.Fatal(err)
	}
	l, err := log.NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	port, _ := strconv.Atoi(os.Getenv("TEST_POSTGRES_PORT"))
	ctx, done := context.WithTimeout(context.Background(), time.Second*10)
	defer done()
	db, err := pkgpg.NewDB(ctx, pkgpg.Config{
		Address:      os.Getenv("TEST_POSTGRES_HOST"),
		Port:         uint(port),
		User:         os.Getenv("TEST_POSTGRES_USER"),
		Password:     os.Getenv("TEST_POSTGRES_PASSWORD"),
		DatabaseName: os.Getenv("TEST_POSTGRES_DB"),
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := apg.NewRepository(db, l)
	if err != nil {
		_ = db.Close()
		t.Fatal(err)
	}
	return out
}

func newTestRefTypeManager(t *testing.T) *api.RefTypeManager {
	repo := newPgRepo(t)
	out, err := api.NewRefTypeManager(api.RefTypeConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
