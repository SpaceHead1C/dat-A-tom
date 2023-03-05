package test

import (
	"context"
	"datatom/internal/domain"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestSetStoredConfig(t *testing.T) {
	mngr := newTestStoredConfigsManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := mngr.Set(
		ctx,
		domain.StoredConfigTomID,
		uuid.MustParse("12345678-1234-1234-1234-123456789012"),
	); err != nil {
		t.Fatal(err)
	}
}

func TestGetStoredConfig(t *testing.T) {
	mngr := newTestStoredConfigsManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	value, err := mngr.Get(ctx, domain.StoredConfigTomID)
	if err != nil {
		t.Fatal(err)
	}
	var id uuid.UUID
	if err := value.ScanStoredConfigValue(&id); err != nil {
		t.Fatal(err)
	}
	t.Log(id.String())
}
