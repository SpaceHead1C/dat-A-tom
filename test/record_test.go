package test

import (
	"context"
	"datatom/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddRecord(t *testing.T) {
	mngr := newTestRecordManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	id, err := mngr.Add(ctx, domain.AddRecordRequest{
		Name: "Запись",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id.String())
}

func TestUpdateRecord(t *testing.T) {
	mngr := newTestRecordManager(t)
	dm := true
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	o, err := mngr.Update(ctx, domain.UpdRecordRequest{
		ID:           uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		DeletionMark: &dm,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("=== Record ===")
	t.Log("ID:", o.ID.String())
	t.Log("reference type ID:", o.ReferenceTypeID.String())
	t.Log("name:", o.Name)
	t.Log("description:", o.Description)
	t.Log("deletion mark:", o.DeletionMark)
	t.Log("hash sum:", o.Sum)
	t.Log("change at:", o.ChangeAt)
}

func TestGetRecord(t *testing.T) {
	mngr := newTestRecordManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	o, err := mngr.Get(ctx, uuid.MustParse("12345678-1234-1234-1234-123456789012"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("=== Record ===")
	t.Log("ID:", o.ID.String())
	t.Log("reference type ID:", o.ReferenceTypeID.String())
	t.Log("name:", o.Name)
	t.Log("description:", o.Description)
	t.Log("deletion mark:", o.DeletionMark)
	t.Log("hash sum:", o.Sum)
	t.Log("change at:", o.ChangeAt)
}
