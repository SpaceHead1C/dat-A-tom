package test

import (
	"context"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"

	"datatom/internal/adapter/pg"
	"datatom/internal/api"
	pkgpg "datatom/pkg/db/pg"
	"datatom/pkg/log"
	"datatom/test/mocks"

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

func newTestRefTypeMockedManager(t *testing.T) (*api.RefTypeManager, *mocks.RefTypeRepository, *mocks.RefTypeBroker) {
	repo := mocks.NewRefTypeRepository(t)
	broker := mocks.NewRefTypeBroker(t)
	out, err := api.NewRefTypeManager(api.RefTypeConfig{
		Repository: repo,
		Broker:     broker,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out, repo, broker
}

func newTestRecordMockedManager(t *testing.T) (*api.RecordManager, *mocks.RecordRepository, *mocks.RecordBroker) {
	repo := mocks.NewRecordRepository(t)
	broker := mocks.NewRecordBroker(t)
	out, err := api.NewRecordManager(api.RecordConfig{
		Repository: repo,
		Broker:     broker,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out, repo, broker
}

func newTestPropertyMockedManager(t *testing.T) (*api.PropertyManager, *mocks.PropertyRepository, *mocks.PropertyBroker) {
	repo := mocks.NewPropertyRepository(t)
	broker := mocks.NewPropertyBroker(t)
	out, err := api.NewPropertyManager(api.PropertyConfig{
		Repository: repo,
		Broker:     broker,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out, repo, broker
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

func newTestChangedDataManager(t *testing.T) *api.ChangedDataManager {
	repo := newPgRepo(t)
	out, err := api.NewChangedDataManager(api.ChangedDataConfig{
		Repository: repo,
		Timeout:    time.Second * 5,
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

func funcName(t *testing.T, f any) string {
	if reflect.ValueOf(f).Kind() != reflect.Func {
		t.Fatalf("%v is not a function", f)
	}
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
