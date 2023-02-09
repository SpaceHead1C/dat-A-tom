package test

import (
	"context"
	. "datatom/internal/domain"
	"testing"
	"time"
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
