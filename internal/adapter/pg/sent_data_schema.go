package pg

import (
	. "datatom/internal/domain"

	"github.com/google/uuid"
	"time"
)

type RecordSentStateSchema struct {
	ID     uuid.UUID `json:"id"`
	Sum    string    `json:"sum"`
	SentAt time.Time `json:"sent_at"`
}

func (rss RecordSentStateSchema) SentData() *RecordSentState {
	return &RecordSentState{
		ID:     rss.ID,
		Sum:    rss.Sum,
		SentAt: rss.SentAt,
	}
}

type PropertySentStateSchema struct {
	ID     uuid.UUID `json:"id"`
	Sum    string    `json:"sum"`
	SentAt time.Time `json:"sent_at"`
}

func (pss PropertySentStateSchema) SentData() *PropertySentState {
	return &PropertySentState{
		ID:     pss.ID,
		Sum:    pss.Sum,
		SentAt: pss.SentAt,
	}
}

type RefTypeSentStateSchema struct {
	ID     uuid.UUID `json:"id"`
	Sum    string    `json:"sum"`
	SentAt time.Time `json:"sent_at"`
}

func (rtss RefTypeSentStateSchema) SentData() *RefTypeSentState {
	return &RefTypeSentState{
		ID:     rtss.ID,
		Sum:    rtss.Sum,
		SentAt: rtss.SentAt,
	}
}

type ValueSentStateSchema struct {
	RecordID   uuid.UUID `json:"record_id"`
	PropertyID uuid.UUID `json:"property_id"`
	Sum        string    `json:"sum"`
	SentAt     time.Time `json:"sent_at"`
}

func (vss ValueSentStateSchema) SentData() *ValueSentState {
	return &ValueSentState{
		RecordID:   vss.RecordID,
		PropertyID: vss.PropertyID,
		Sum:        vss.Sum,
		SentAt:     vss.SentAt,
	}
}
