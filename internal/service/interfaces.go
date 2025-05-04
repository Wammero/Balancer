package service

import (
	"context"

	"github.com/Wammero/Balancer/internal/models"
)

type ClientService interface {
	CreateClient(ctx context.Context, client_id string, capacity, ratePerSec int) (string, error)
	ListClients(ctx context.Context) (*[]models.TokenBucket, error)
	GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error)
	DeleteClient(ctx context.Context, clientID string) error
	UpdateClient(ctx context.Context, clientID string, capacity, ratePerSec int) error
}

type ProxyService interface {
	CheckRateLimit(ctx context.Context) error
}
