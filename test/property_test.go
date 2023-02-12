package test

import (
	"context"
	"datatom/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddProperty(t *testing.T) {
	mngr := newTestPropertyManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20000)
	defer cancel()
	id, err := mngr.Add(ctx, domain.AddPropertyRequest{
		Name:           "Запись",
		Types:          []domain.Type{domain.TypeText, domain.TypeReference},
		RefTypeIDs:     []uuid.UUID{uuid.MustParse("12345678-1234-1234-1234-123456789012")},
		OwnerRefTypeID: uuid.Nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id.String())
}

func TestUpdateProperty(t *testing.T) {
	mngr := newTestPropertyManager(t)
	description := "dscr"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	o, err := mngr.Update(ctx, domain.UpdPropertyRequest{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Description: &description,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("=== Property ===")
	t.Log("ID:", o.ID.String())
	t.Log("name:", o.Name)
	t.Log("description:", o.Description)
	t.Log("types:", o.Types)
	t.Log("reference types:", o.RefTypeIDs)
	t.Log("owner type ID:", o.OwnerRefTypeID.String())
	t.Log("hash sum:", o.Sum)
	t.Log("change at:", o.ChangeAt)
}

func TestGetProperty(t *testing.T) {
	mngr := newTestPropertyManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	o, err := mngr.Get(ctx, uuid.MustParse("12345678-1234-1234-1234-123456789012"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("=== Property ===")
	t.Log("ID:", o.ID.String())
	t.Log("name:", o.Name)
	t.Log("description:", o.Description)
	t.Log("types:", o.Types)
	t.Log("reference types:", o.RefTypeIDs)
	t.Log("owner type ID:", o.OwnerRefTypeID.String())
	t.Log("hash sum:", o.Sum)
	t.Log("change at:", o.ChangeAt)
}
