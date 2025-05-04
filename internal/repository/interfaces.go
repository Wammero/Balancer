package repository

import (
	"context"

	"github.com/Wammero/Balancer/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ClientRepository interface {
	Pool() *pgxpool.Pool
	CreateClient(ctx context.Context, tx pgx.Tx, cleint_id string, capacity, ratePerSecond int) error
	GetClients(ctx context.Context) (*[]models.TokenBucket, error)
	GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error)
	DeleteClient(ctx context.Context, tx pgx.Tx, clientID string) error
	UpdateClient(ctx context.Context, tx pgx.Tx, clientID string, capacity, ratePerSecond int) error
}

type ProxyRepository interface {
	Pool() *pgxpool.Pool
	GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error)
	UpdateTokens(ctx context.Context, tx pgx.Tx, clientID string, tokens int) (*models.TokenBucket, error)
}
