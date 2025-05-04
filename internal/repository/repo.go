package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	ClientRepository ClientRepository
	ProxyRepository  ProxyRepository
	pool             *pgxpool.Pool
}

func New(connstr string) (*Repository, error) {
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		ClientRepository: NewClientRepository(pool),
		ProxyRepository:  NewProxyRepository(pool),
		pool:             pool,
	}, nil
}

func (r *Repository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *Repository) Pool() *pgxpool.Pool {
	return r.pool
}
