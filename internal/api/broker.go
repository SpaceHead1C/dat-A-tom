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

type RecordSender struct {
	man *RecordManager
	req domain.SendRecordRequest
}

func (rs *RecordSender) Send(ctx context.Context) error {
	return rs.man.Send(ctx, rs.req)
}

func (rs *RecordSender) SumEqualsSent(ctx context.Context, transaction db.Transaction) (bool, error) {
	state, err := rs.man.GetSentState(ctx, rs.req.ID, transaction)
	if err != nil {
		if errors.Is(err, domain.ErrSentDataNotFound) {
			return false, nil
		}
		return false, err
	}
	return state.Sum == rs.req.Sum, nil
}

func (rs *RecordSender) SetSentState(ctx context.Context, transaction db.Transaction) error {
	_, err := rs.man.SetSentState(
		ctx,
		domain.RecordSentState{
			ID:     rs.req.ID,
			Sum:    rs.req.Sum,
			SentAt: time.Now().UTC(),
		},
		transaction,
	)
	return err
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

type PropertySender struct {
	man *PropertyManager
	req domain.SendPropertyRequest
}

func (ps *PropertySender) Send(ctx context.Context) error {
	return ps.man.Send(ctx, ps.req)
}

func (ps *PropertySender) SumEqualsSent(ctx context.Context, transaction db.Transaction) (bool, error) {
	state, err := ps.man.GetSentState(ctx, ps.req.ID, transaction)
	if err != nil {
		if errors.Is(err, domain.ErrSentDataNotFound) {
			return false, nil
		}
		return false, err
	}
	return state.Sum == ps.req.Sum, nil
}

func (ps *PropertySender) SetSentState(ctx context.Context, transaction db.Transaction) error {
	_, err := ps.man.SetSentState(
		ctx,
		domain.PropertySentState{
			ID:     ps.req.ID,
			Sum:    ps.req.Sum,
			SentAt: time.Now().UTC(),
		},
		transaction,
	)
	return err
}
