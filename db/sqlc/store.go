package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type SQLStore struct {
	db *pgxpool.Pool
	*Queries
}

func NewStore(dbSource string) (Store, error) {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	store := &SQLStore{
		db:      connPool,
		Queries: New(connPool),
	}

	return store, nil
}

func (store *SQLStore) GetDB() *pgxpool.Pool {
	return store.db
}
