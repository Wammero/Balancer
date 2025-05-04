package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wammero/Balancer/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type clientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) *clientRepository {
	return &clientRepository{pool: pool}
}

func (c *clientRepository) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *clientRepository) CreateClient(ctx context.Context, tx pgx.Tx, clientID string, capacity, ratePerSecond int) error {
	query := `
        INSERT INTO token_buckets (key, tokens, last_refill, tokens_per_second, capacity)
        VALUES ($1, $2, NOW(), $3, $4)
    `
	_, err := tx.Exec(ctx, query, clientID, capacity, ratePerSecond, capacity)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("клиент с ID %s уже существует", clientID)
		}
		return fmt.Errorf("не удалось создать клиента: %w", err)
	}
	return nil
}

func (c *clientRepository) GetClients(ctx context.Context) (*[]models.TokenBucket, error) {
	query := `
		SELECT key, tokens, last_refill, tokens_per_second, capacity, created_at, updated_at
		FROM token_buckets
	`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить клиентов: %w", err)
	}
	defer rows.Close()

	var clients []models.TokenBucket
	for rows.Next() {
		var client models.TokenBucket
		err := rows.Scan(
			&client.Key,
			&client.Tokens,
			&client.LastRefill,
			&client.TokensPerSecond,
			&client.Capacity,
			&client.CreatedAt,
			&client.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении строки: %w", err)
		}
		clients = append(clients, client)
	}

	return &clients, nil
}

func (c *clientRepository) GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error) {
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

func (c *clientRepository) DeleteClient(ctx context.Context, tx pgx.Tx, clientID string) error {
	query := `DELETE FROM token_buckets WHERE key = $1`

	cmdTag, err := c.pool.Exec(ctx, query, clientID)
	if err != nil {
		return fmt.Errorf("не удалось удалить клиента: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("клиент с id %s не найден", clientID)
	}

	return nil
}

func (c *clientRepository) UpdateClient(ctx context.Context, tx pgx.Tx, clientID string, capacity, ratePerSecond int) error {
	query := `
		UPDATE token_buckets
		SET tokens_per_second = $1,
		    capacity = $2,
			tokens = LEAST($2, tokens),
		    updated_at = now()
		WHERE key = $3
	`

	cmdTag, err := c.pool.Exec(ctx, query, ratePerSecond, capacity, clientID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных клиента: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("клиент с id %s не найден", clientID)
	}

	return nil
}
