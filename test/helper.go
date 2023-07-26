package test

import (
	"context"
	"datatom/internal/adapter/pg"
	"datatom/internal/api"
	pkgpg "datatom/pkg/db/pg"
	"datatom/pkg/log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/subosito/gotenv"
)

func newPgRepo(t *testing.T) *pg.Repository {
	if err := gotenv.Load(); err != nil {
		t.Fatal(err)
	}
	l, err := log.NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(os.Getenv("TEST_POSTGRES_PORT"))
	if err != nil {
		t.Fatal(err)
	}
	db, err := pkgpg.NewPoolConfig(pkgpg.Config{
		Address:      os.Getenv("TEST_POSTGRES_HOST"),
		Port:         uint(port),
		User:         os.Getenv("TEST_POSTGRES_USER"),
		Password:     os.Getenv("TEST_POSTGRES_PASSWORD"),
		DatabaseName: os.Getenv("TEST_POSTGRES_DB"),
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	out, err := pg.NewRepository(ctx, pg.Config{
		ConnectConfig: db,
		Logger:        l,
	})
	cancel()
	if err != nil {
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

func newTestRecordManager(t *testing.T) *api.RecordManager {
	repo := newPgRepo(t)
	out, err := api.NewRecordManager(api.RecordConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func newTestPropertyManager(t *testing.T) *api.PropertyManager {
	repo := newPgRepo(t)
	out, err := api.NewPropertyManager(api.PropertyConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func newTestValueManager(t *testing.T) *api.ValueManager {
	repo := newPgRepo(t)
	out, err := api.NewValueManager(api.ValueConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func newTestSentDataManager(t *testing.T) *api.SentDataManager {
	repo := newPgRepo(t)
	out, err := api.NewSentDataManager(api.SentDataConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func newTestStoredConfigsManager(t *testing.T) *api.StoredConfigsManager {
	repo := newPgRepo(t)
	out, err := api.NewStoredConfigManager(api.StoredConfigsConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
