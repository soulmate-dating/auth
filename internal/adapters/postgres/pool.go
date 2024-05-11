package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TxCtxKey struct{}

type Connection interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type ConnPool interface {
	GetTx(ctx context.Context) Connection
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type Pool struct {
	pool Connection
}

func NewPool(pool Connection) *Pool {
	return &Pool{pool: pool}
}

func (p *Pool) GetTx(ctx context.Context) Connection {
	if tx := p.AcquireTx(ctx); tx != nil {
		return tx
	}
	return p.pool
}

func (p *Pool) AcquireTx(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(TxCtxKey{}).(pgx.Tx)

	return tx
}

func (p *Pool) RunInTx(ctx context.Context, f func(context.Context) error) error {
	if tx := p.AcquireTx(ctx); tx != nil {
		return f(ctx)
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
	}()

	if err := f(context.WithValue(ctx, TxCtxKey{}, tx)); err != nil {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil {
			err = fmt.Errorf("%w - rollback transaction: %w", err, errRollback)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
