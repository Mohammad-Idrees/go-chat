package db

import (
	"context"
	"fmt"
	"log"
	"project/models"
)

type Repository interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

type RepositoryImpl struct {
	*Queries
	db *models.Database
}

func NewRepository(db *models.Database) Repository {
	return &RepositoryImpl{
		Queries: New(db.ConnPool),
		db:      db,
	}
}

// ExecTx implements Repository.
func (r *RepositoryImpl) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := r.db.ConnPool.Begin(ctx)
	if err != nil {
		log.Println("failed to begin tx", err)
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)

}
