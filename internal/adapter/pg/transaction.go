package pg

import (
	"context"
	"datatom/internal/domain"
	"datatom/pkg/db"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	pgx.Tx
}

func (r *Repository) BeginTransaction(ctx context.Context) (db.Transaction, error) {
	tx, err := r.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	return Transaction{tx}, nil
}

func (t Transaction) Begin(ctx context.Context) (db.Transaction, error) {
	out, err := t.Tx.Begin(ctx)
	return Transaction{out}, err
}

func (t Transaction) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}

func (t Transaction) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func unwrapTransaction(t db.Transaction) (pgx.Tx, error) {
	switch t.(type) {
	case Transaction:
		return t.(Transaction).Tx, nil
	default:
		return nil, fmt.Errorf("%w %T of transaction", domain.ErrUnexpectedType, t)
	}
}

func funcQueryRow(r *Repository, t db.Transaction) (func(context.Context, string, ...any) pgx.Row, error) {
	if t == nil {
		return r.QueryRow, nil
	}
	tx, err := unwrapTransaction(t)
	if err != nil {
		return nil, err
	}
	return tx.QueryRow, nil
}
