package domain

import (
	"context"

	"github.com/google/uuid"
)

type RefTypeRepository interface {
	AddRefType(context.Context, AddRefTypeRequest) (uuid.UUID, error)
}

type AddRefTypeRequest struct {
	Name        string
	Description string
}
