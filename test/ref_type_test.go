package test

import (
	. "datatom/internal/domain"
	"testing"
)

func TestAddRefType(t *testing.T) {
	mngr := newTestRefTypeManager(t)
	id, err := mngr.Add(AddRefTypeRequest{
		Name: "Группы магазинов для графиков",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id.String())
}
