package domain

import (
	"context"
	"datatom/pkg/db"
	"time"

	"github.com/google/uuid"
)

const DeliveryTypeRefType = "reference_type"

type RefTypeRepository interface {
	AddRefType(context.Context, AddRefTypeRequest) (uuid.UUID, error)
	UpdateRefType(context.Context, UpdRefTypeRequest) (*RefType, error)
	GetRefType(context.Context, uuid.UUID) (*RefType, error)
	GetRefTypeByKey(context.Context, []byte) (*RefType, error)
	GetRefTypeSentStateForUpdate(context.Context, uuid.UUID, db.Transaction) (*RefTypeSentState, error)
	SetSentRefType(context.Context, RefTypeSentState, db.Transaction) (*RefTypeSentState, error)
}

type RefTypeBroker interface {
	SendRefType(context.Context, SendRefTypeRequest) error
}

type RefType struct {
	ID          uuid.UUID
	Name        string
	Description string
	Sum         string
	ChangeAt    time.Time
}

type RefTypeSentState struct {
	ID     uuid.UUID
	Sum    string
	SentAt time.Time
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

type SendRefTypeRequest struct {
	RefType
	TomID       uuid.UUID
	Exchange    string
	RoutingKeys []string
}
