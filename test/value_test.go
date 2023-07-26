package test

import (
	"context"
	"datatom/internal/domain"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSetValue(t *testing.T) {
	mngr := newTestValueManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	o, err := mngr.Set(ctx, domain.SetValueRequest{
		RecordID:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		PropertyID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Type:       domain.TypeText,
		Value:      "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("=== Value ===")
	t.Log("record ID:", o.RecordID.String())
	t.Log("property ID:", o.PropertyID.String())
	t.Log("type:", o.Type.String())
	t.Log("reference type ID:", o.RefTypeID.String())
	t.Log("value:", fmt.Sprintf("%v", o.Value))
	t.Log("hash sum:", o.Sum)
	t.Log("change at:", o.ChangeAt)
}

func TestChangedValues(t *testing.T) {
	mngr := newTestValueManager(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	vs, err := mngr.ChangedValues(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range vs {
		t.Logf(`=== Value ===
record ID: %s
property ID: %s
type: %s
reference type ID: %s
value: %v
hash sum: %s
change at: %s`,
			v.RecordID.String(), v.PropertyID.String(), v.Type.String(), v.RefTypeID.String(), v.Value, v.Sum, v.ChangeAt,
		)
	}
}
