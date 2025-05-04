package service

import (
	"context"
	"fmt"

	"github.com/Wammero/Balancer/internal/cache"
	"github.com/Wammero/Balancer/internal/models"
	"github.com/Wammero/Balancer/internal/repository"
	"github.com/Wammero/Balancer/pkg/jwt"
)

type clientService struct {
	repo  repository.ClientRepository
	cache cache.ClientCache
}

func NewClientService(repo repository.ClientRepository, cache cache.ClientCache) *clientService {
	return &clientService{repo: repo, cache: cache}
}

func (c *clientService) CreateClient(ctx context.Context, clientID string, capacity, ratePerSec int) (token string, err error) {
	tx, err := c.repo.Pool().Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	err = c.repo.CreateClient(ctx, tx, clientID, capacity, ratePerSec)
	if err != nil {
		return "", err
	}

	err = c.cache.CreateClient(ctx, clientID, capacity, ratePerSec)
	if err != nil {
		return "", fmt.Errorf("не удалось обновить кэш для клиента %s: %w", clientID, err)
	}

	token, err = jwt.GenerateJWT(clientID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *clientService) ListClients(ctx context.Context) (*[]models.TokenBucket, error) {
	clients, err := c.repo.GetClients(ctx)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (c *clientService) GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error) {
	client, err := c.cache.GetClientByID(ctx, clientID)
	if err == nil {
		return client, nil
	}

	client, err = c.repo.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if err := c.cache.AddClient(ctx, client); err != nil {
		return nil, fmt.Errorf("не удалось обновить кэш для клиента %s: %w", clientID, err)
	}

	return client, nil
}

func (c *clientService) DeleteClient(ctx context.Context, clientID string) error {
	tx, err := c.repo.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	err = c.repo.DeleteClient(ctx, tx, clientID)
	if err != nil {
		return err
	}

	err = c.cache.DeleteClient(ctx, clientID)
	if err != nil {
		return fmt.Errorf("не удалось удалить кэш для клиента %s: %w", clientID, err)
	}

	return nil
}

func (c *clientService) UpdateClient(ctx context.Context, clientID string, capacity, ratePerSec int) error {
	tx, err := c.repo.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	err = c.repo.UpdateClient(ctx, tx, clientID, capacity, ratePerSec)
	if err != nil {
		return err
	}

	err = c.cache.UpdateClient(ctx, clientID, capacity, ratePerSec)
	if err != nil {
		return fmt.Errorf("не удалось обновить кэш для клиента %s: %w", clientID, err)
	}

	return nil
}
