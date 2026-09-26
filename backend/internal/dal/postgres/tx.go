package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var ErrNoTx = errors.New("no transactions found")

type ctxValueKey string

const txCtxValueKey ctxValueKey = "transaction"

func (p *Postgres) Beginx(ctx context.Context) (context.Context, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		p.log.Error(
			"failed to begin transaction",
			"error", err,
		)

		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return context.WithValue(ctx, txCtxValueKey, tx), nil
}

func (p *Postgres) BeginTxx(ctx context.Context, opts *sql.TxOptions) (context.Context, error) {
	tx, err := p.db.BeginTxx(ctx, opts)
	if err != nil {
		p.log.Error(
			"failed to begin transaction",
			"error", err,
		)

		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return context.WithValue(ctx, txCtxValueKey, tx), nil
}

func (p *Postgres) Rollback(ctx context.Context) {
	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if !ok {
		p.log.Error(
			"nothing to rollback",
			"error", ErrNoTx,
		)

		return
	}

	if err := tx.Rollback(); err != nil {
		p.log.Error(
			"failed to rollback transaction",
			"error", err,
		)
	}
}

func (p *Postgres) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(txCtxValueKey).(*sqlx.Tx)
	if !ok {
		p.log.Error(
			"nothing to commit",
			"error", ErrNoTx,
		)

		return ErrNoTx
	}

	if err := tx.Commit(); err != nil {
		p.log.Error(
			"failed to commit transaction",
			"error", err,
		)

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
