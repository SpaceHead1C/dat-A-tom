package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type SentDataRepository interface {
	GetSentData(context.Context, GetSentDataRequest) (*SentData, error)
}

type SentData struct {
	RecordID   uuid.UUID
	PropertyID uuid.UUID
	Sum        string
	SentAt     time.Time
}

type GetSentDataRequest struct {
	RecordID   uuid.UUID
	PropertyID uuid.UUID
}
