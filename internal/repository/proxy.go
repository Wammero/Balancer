package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wammero/Balancer/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type proxyRepository struct {
	pool *pgxpool.Pool
}

func NewProxyRepository(pool *pgxpool.Pool) *proxyRepository {
	return &proxyRepository{pool: pool}
}

func (p *proxyRepository) Pool() *pgxpool.Pool {
	return p.pool
}

func (c *proxyRepository) GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error) {
	query := `
		SELECT key, tokens, last_refill, tokens_per_second, capacity, created_at, updated_at
		FROM token_buckets
		WHERE key = $1
	`

	var client models.TokenBucket
	err := c.pool.QueryRow(ctx, query, clientID).Scan(
		&client.Key,
		&client.Tokens,
		&client.LastRefill,
		&client.TokensPerSecond,
		&client.Capacity,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("клиент с id %s не найден", clientID)
		}
		return nil, fmt.Errorf("ошибка при получении клиента: %w", err)
	}

	return &client, nil
}

func (c *proxyRepository) UpdateTokens(ctx context.Context, tx pgx.Tx, clientID string, tokens int) (*models.TokenBucket, error) {
	query := `
		UPDATE token_buckets
		SET tokens = $1,
			last_refill = now(),
		    updated_at = now()
		WHERE key = $2
		RETURNING key, tokens, last_refill, tokens_per_second, capacity, created_at, updated_at
	`

	row := c.pool.QueryRow(ctx, query, tokens, clientID)

	var bucket models.TokenBucket
	err := row.Scan(
		&bucket.Key,
		&bucket.Tokens,
		&bucket.LastRefill,
		&bucket.TokensPerSecond,
		&bucket.Capacity,
		&bucket.CreatedAt,
		&bucket.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("клиент с id %s не найден", clientID)
		}
		return nil, fmt.Errorf("ошибка при обновлении данных клиента: %w", err)
	}

	return &bucket, nil
}
