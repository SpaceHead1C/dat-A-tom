package api

import (
	"context"
	"datatom/internal/domain"
	"datatom/pkg/db"
	"errors"
	"time"
)

type Sender interface {
	Send(context.Context) error
	SumEqualsSent(context.Context, db.Transaction) (bool, error)
	SetSentState(context.Context, db.Transaction) error
}

type RefTypeSender struct {
	man *RefTypeManager
	req domain.SendRefTypeRequest
}

func (rts *RefTypeSender) Send(ctx context.Context) error {
	return rts.man.Send(ctx, rts.req)
}

func (rts *RefTypeSender) SumEqualsSent(ctx context.Context, transaction db.Transaction) (bool, error) {
	state, err := rts.man.GetSentState(ctx, rts.req.ID, transaction)
	if err != nil {
		if errors.Is(err, domain.ErrSentDataNotFound) {
			return false, nil
		}
		return false, err
	}
	return state.Sum == rts.req.Sum, nil
}

func (rts *RefTypeSender) SetSentState(ctx context.Context, transaction db.Transaction) error {
	_, err := rts.man.SetSentState(
		ctx,
		domain.RefTypeSentState{
			ID:     rts.req.ID,
			Sum:    rts.req.Sum,
			SentAt: time.Now().UTC(),
		},
		transaction,
	)
	return err
}
