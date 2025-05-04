package service

import (
	"context"
	"fmt"

	"github.com/Wammero/Balancer/internal/cache"
	"github.com/Wammero/Balancer/internal/limiter"
	"github.com/Wammero/Balancer/internal/models"
	"github.com/Wammero/Balancer/internal/repository"
	"github.com/Wammero/Balancer/pkg/jwt"
)

type proxyService struct {
	repo        repository.ProxyRepository
	cache       cache.ProxyCache
	rateLimiter limiter.Limiter
}

func NewProxyService(repo repository.ProxyRepository, cache cache.ProxyCache, limiter limiter.Limiter) *proxyService {
	return &proxyService{
		repo:        repo,
		cache:       cache,
		rateLimiter: limiter,
	}
}

func (p *proxyService) CheckRateLimit(ctx context.Context) error {
	clientID, ok := jwt.GetUserID(ctx)
	if !ok {
		return fmt.Errorf("не удалось получить userID из JWT токена")
	}

	var client *models.TokenBucket
	var err error

	client, err = p.cache.GetClientByID(ctx, clientID)
	if err != nil {
		client, err = p.repo.GetClientByID(ctx, clientID)
		if err != nil {
			return fmt.Errorf("ошибка получения клиента из базы: %w", err)
		}
		if err := p.cache.AddClient(ctx, client); err != nil {
			return fmt.Errorf("не удалось обновить кэш для клиента %s: %w", clientID, err)
		}
	}

	if err := p.rateLimiter.Check(ctx, client); err != nil {
		return err
	}

	tx, err := p.repo.Pool().Begin(ctx)
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

	client, err = p.repo.UpdateTokens(ctx, tx, clientID, client.Tokens)
	if err != nil {
		return fmt.Errorf("не удалось обновить токены: %w", err)
	}

	err = p.cache.AddClient(ctx, client)
	if err != nil {
		return fmt.Errorf("не удалось сохранить состояние клиента в кэше: %w", err)
	}

	return nil

}
