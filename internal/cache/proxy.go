package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Wammero/Balancer/internal/models"
	"github.com/redis/go-redis/v9"
)

type proxyCache struct {
	client *redis.Client
}

func NewProxyCache(client *redis.Client) *proxyCache {
	return &proxyCache{client: client}
}

func (c *proxyCache) GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error) {
	var client models.TokenBucket

	data, err := c.client.Get(ctx, clientID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("клиент с id %s не найден", clientID)
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось получить клиента из Redis: %w", err)
	}

	err = json.Unmarshal([]byte(data), &client)
	if err != nil {
		return nil, fmt.Errorf("не удалось десериализовать данные клиента: %w", err)
	}

	return &client, nil
}

func (c *proxyCache) AddClient(ctx context.Context, client *models.TokenBucket) error {
	clientData, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать данные клиента: %w", err)
	}

	err = c.client.Set(ctx, client.Key, clientData, 0).Err()
	if err != nil {
		return fmt.Errorf("не удалось добавить клиента в кэш Redis: %w", err)
	}

	return nil
}
