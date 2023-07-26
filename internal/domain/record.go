package domain

import (
	"context"
	"datatom/pkg/db"
	"time"

	"github.com/google/uuid"
)

const DeliveryTypeRecord = "record"

type RecordRepository interface {
	AddRecord(context.Context, AddRecordRequest) (uuid.UUID, error)
	UpdateRecord(context.Context, UpdRecordRequest) (*Record, error)
	GetRecord(context.Context, uuid.UUID) (*Record, error)
	GetRecordByKey(context.Context, []byte) (*Record, error)
	GetRecordSentStateForUpdate(context.Context, uuid.UUID, db.Transaction) (*RecordSentState, error)
	SetSentRecord(context.Context, RecordSentState, db.Transaction) (*RecordSentState, error)
}

type RecordBroker interface {
	SendRecord(context.Context, SendRecordRequest) error
}

type Record struct {
	ID              uuid.UUID
	ReferenceTypeID uuid.UUID
	Name            string
	Description     string
	DeletionMark    bool
	Sum             string
	ChangeAt        time.Time
}

type RecordSentState struct {
	ID     uuid.UUID
	Sum    string
	SentAt time.Time
}

type AddRecordRequest struct {
	Name            string
	Description     string
	DeletionMark    bool
	ReferenceTypeID uuid.UUID
}

type UpdRecordRequest struct {
	ID           uuid.UUID
	Name         *string
	Description  *string
	DeletionMark *bool
}

type SendRecordRequest struct {
	Record
	TomID       uuid.UUID
	Exchange    string
	RoutingKeys []string
}
