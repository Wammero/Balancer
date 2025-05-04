package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Wammero/Balancer/internal/models"
	"github.com/redis/go-redis/v9"
)

type clientCache struct {
	client *redis.Client
}

func NewClientCache(client *redis.Client) *clientCache {
	return &clientCache{client: client}
}

func (c *clientCache) CreateClient(ctx context.Context, clientID string, capacity, ratePerSecond int) error {
	client := models.TokenBucket{
		Key:             clientID,
		Tokens:          capacity,
		LastRefill:      time.Now(),
		TokensPerSecond: ratePerSecond,
		Capacity:        capacity,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	data, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать клиента: %w", err)
	}

	err = c.client.Set(ctx, clientID, data, 0).Err()
	if err != nil {
		return fmt.Errorf("не удалось сохранить клиента в Redis: %w", err)
	}
	return nil
}

func (c *clientCache) GetClientByID(ctx context.Context, clientID string) (*models.TokenBucket, error) {
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

func (c *clientCache) UpdateClient(ctx context.Context, clientID string, capacity, ratePerSecond int) error {
	client, err := c.GetClientByID(ctx, clientID)
	if err != nil {
		return err
	}

	client.Capacity = capacity
	client.TokensPerSecond = ratePerSecond
	client.Tokens = min(capacity, client.Tokens)
	client.UpdatedAt = time.Now()

	data, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать обновленного клиента: %w", err)
	}

	err = c.client.Set(ctx, clientID, data, 0).Err()
	if err != nil {
		return fmt.Errorf("не удалось обновить клиента в Redis: %w", err)
	}

	return nil
}

func (c *clientCache) DeleteClient(ctx context.Context, clientID string) error {
	err := c.client.Del(ctx, clientID).Err()
	if err != nil {
		return fmt.Errorf("не удалось удалить клиента из Redis: %w", err)
	}
	return nil
}

func (c *clientCache) AddClient(ctx context.Context, client *models.TokenBucket) error {
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
