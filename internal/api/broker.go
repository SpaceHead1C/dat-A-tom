package api

import (
	"context"
	"datatom/pkg/db"
)

type Sender interface {
	Send(context.Context) error
	SumEqualsSent(context.Context, db.Transaction) (bool, error)
	SetSentState(context.Context, db.Transaction) error
}
