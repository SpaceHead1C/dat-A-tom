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
		Name:         "Запись",
		Types:        []domain.Type{domain.TypeText, domain.TypeReference},
		RefTypes:     []uuid.UUID{uuid.MustParse("12345678-1234-1234-1234-123456789012")},
		OwnerRefType: uuid.Nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id.String())
}
