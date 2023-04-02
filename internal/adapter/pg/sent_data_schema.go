package pg

import (
	. "datatom/internal/domain"

	"github.com/google/uuid"
	"time"
)

type SentDataSchema struct {
	RecordID   uuid.UUID `json:"record_id"`
	PropertyID uuid.UUID `json:"property_id"`
	Sum        string    `json:"sum"`
	SentAt     time.Time `json:"sent_at"`
}

func (sds *SentDataSchema) SentData() *SentData {
	return &SentData{
		RecordID:   sds.RecordID,
		PropertyID: sds.PropertyID,
		Sum:        sds.Sum,
		SentAt:     sds.SentAt,
	}
}
