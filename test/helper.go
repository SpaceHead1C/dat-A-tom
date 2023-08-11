package test

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"datatom/internal/api"
	"datatom/test/mocks"
)

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

func newTestValueMockedManager(t *testing.T) (*api.ValueManager, *mocks.ValueRepository, *mocks.ValueBroker) {
	repo := mocks.NewValueRepository(t)
	broker := mocks.NewValueBroker(t)
	out, err := api.NewValueManager(api.ValueConfig{
		Repository: repo,
		Broker:     broker,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out, repo, broker
}

func newTestChangedDataManager(t *testing.T) (*api.ChangedDataManager, *mocks.ChangedDataRepository) {
	repo := mocks.NewChangedDataRepository(t)
	out, err := api.NewChangedDataManager(api.ChangedDataConfig{
		Repository: repo,
		Timeout:    time.Second * 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out, repo
}

func newTestStoredConfigsManager(t *testing.T) (*api.StoredConfigsManager, *mocks.StoredConfigRepository) {
	repo := mocks.NewStoredConfigRepository(t)
	out, err := api.NewStoredConfigManager(api.StoredConfigsConfig{
		Repository: repo,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out, repo
}

func funcName(t *testing.T, f any) string {
	if reflect.ValueOf(f).Kind() != reflect.Func {
		t.Fatalf("%v is not a function", f)
	}
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
