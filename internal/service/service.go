package service

import (
	"github.com/Wammero/Balancer/internal/cache"
	ratelimiter "github.com/Wammero/Balancer/internal/limiter"
	"github.com/Wammero/Balancer/internal/repository"
)

type Service struct {
	ClientService ClientService
	ProxyService  ProxyService
}

func New(repo *repository.Repository, redis *cache.RedisCache, limiter ratelimiter.Limiter) *Service {
	return &Service{
		ClientService: NewClientService(repo.ClientRepository, redis.ClientCache),
		ProxyService:  NewProxyService(repo.ProxyRepository, redis.ProxyCache, limiter),
	}
}
