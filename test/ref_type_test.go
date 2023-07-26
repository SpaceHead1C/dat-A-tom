package test

import (
	"context"
	. "datatom/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddRefType(t *testing.T) {
	mngr := newTestRefTypeManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	id, err := mngr.Add(ctx, AddRefTypeRequest{
		Name: "Группы магазинов для графиков",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id.String())
}

func TestUpdateRefType(t *testing.T) {
	mngr := newTestRefTypeManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	description := "dscr"
	rt, err := mngr.Update(ctx, UpdRefTypeRequest{
		ID:          uuid.MustParse("12345678-1234-1234-1234-123456789012"),
		Description: &description,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("=== Reference type ===")
	t.Log("ID:", rt.ID.String())
	t.Log("name:", rt.Name)
	t.Log("description:", rt.Description)
}

func TestGetRefType(t *testing.T) {
	mngr := newTestRefTypeManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	rt, err := mngr.Get(ctx, uuid.MustParse("12345678-1234-1234-1234-123456789012"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("=== Reference type ===")
	t.Log("ID:", rt.ID.String())
	t.Log("name:", rt.Name)
	t.Log("description:", rt.Description)
	t.Log("sum:", rt.Sum)
	t.Log("change at:", rt.ChangeAt)
}
