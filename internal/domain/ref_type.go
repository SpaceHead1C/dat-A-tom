package domain

import (
	"context"

	"github.com/google/uuid"
)

const DeliveryTypeRefType = "reference_type"

type RefTypeRepository interface {
	AddRefType(context.Context, AddRefTypeRequest) (uuid.UUID, error)
	UpdateRefType(context.Context, UpdRefTypeRequest) (*RefType, error)
	GetRefType(context.Context, uuid.UUID) (*RefType, error)
}

type RefType struct {
	ID          uuid.UUID
	Name        string
	Description string
}

type AddRefTypeRequest struct {
	Name        string
	Description string
}

type UpdRefTypeRequest struct {
	ID          uuid.UUID
	Name        *string
	Description *string
}
