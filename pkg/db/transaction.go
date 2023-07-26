package db

import (
	"context"
)

type Transaction interface {
	Begin(context.Context) (Transaction, error)
	Rollback(context.Context) error
	Commit(context.Context) error
}

type TransactionBeginner interface {
	BeginTransaction(context.Context) (Transaction, error)
}
