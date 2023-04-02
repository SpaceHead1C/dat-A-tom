package test

import (
	"context"
	"datatom/internal/domain"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestGetSentData(t *testing.T) {
	mngr := newTestSentDataManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	o, err := mngr.Get(ctx, domain.GetSentDataRequest{
		RecordID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		PropertyID: uuid.MustParse("00000002-0000-0000-0000-000000000001"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`=== Sent data ===
record ID: %s
property ID: %s
hash sum: %s
change at: %s`,
		o.RecordID.String(), o.PropertyID.String(), o.Sum, o.SentAt,
	)
}
