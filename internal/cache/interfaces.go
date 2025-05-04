package cache

import (
	"context"

	"github.com/Wammero/Balancer/internal/models"
)

type ClientCache interface {
	CreateClient(ctx context.Context, clientID string, capacity, ratePerSecond int) error
	GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error)
	UpdateClient(ctx context.Context, clientID string, capacity, ratePerSecond int) error
	DeleteClient(ctx context.Context, clientID string) error
	AddClient(ctx context.Context, client *models.TokenBucket) error
}

type ProxyCache interface {
	GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error)
	AddClient(ctx context.Context, client *models.TokenBucket) error
}
