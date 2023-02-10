package test

import (
	"context"
	"datatom/internal/domain"
	"testing"
	"time"
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
