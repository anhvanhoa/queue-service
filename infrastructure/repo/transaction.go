package repo

import (
	"context"
	"errors"

	repository "queue-service/domain/repository"

	"github.com/go-pg/pg/v10"
)

type managerTransaction struct {
	db *pg.DB
}

var ErrTxContextKey error = errors.New("no transaction in context")

func NewManagerTransaction(db *pg.DB) repository.ManagerTransaction {
	return &managerTransaction{
		db: db,
	}
}

type txContextKey struct{}

func (mt *managerTransaction) RunInTransaction(fn func(ctx context.Context) error) error {
	tx, err := mt.db.BeginContext(mt.db.Context())
	if err != nil {
		return err
	}
	txCtx := context.WithValue(mt.db.Context(), txContextKey{}, tx)
	// Nếu có lỗi thì rollback
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}
	// Commit transaction
	return tx.Commit()
}

func (mt *managerTransaction) Begin() (context.Context, error) {
	ctx := mt.db.Context()
	tx, err := mt.db.BeginContext(ctx)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, txContextKey{}, tx), nil
}

func (mt *managerTransaction) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(txContextKey{}).(*pg.Tx)
	if !ok {
		return ErrTxContextKey
	}
	return tx.Commit()
}

func (mt *managerTransaction) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(txContextKey{}).(*pg.Tx)
	if !ok {
		return ErrTxContextKey
	}
	return tx.Rollback()
}
