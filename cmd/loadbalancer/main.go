package main

import (
	"github.com/Wammero/Balancer/internal/balancer"
	"github.com/Wammero/Balancer/internal/cache"
	"github.com/Wammero/Balancer/internal/config"
	"github.com/Wammero/Balancer/internal/handler"
	"github.com/Wammero/Balancer/internal/health"
	"github.com/Wammero/Balancer/internal/limiter"
	"github.com/Wammero/Balancer/internal/logger"
	"github.com/Wammero/Balancer/internal/migration"
	"github.com/Wammero/Balancer/internal/repository"
	"github.com/Wammero/Balancer/internal/router"
	"github.com/Wammero/Balancer/internal/server"
	"github.com/Wammero/Balancer/internal/service"
	"github.com/Wammero/Balancer/pkg/jwt"
)

func main() {
	// Конфигурация
	cfg := config.NewConfig()

	// Инициализация логгера
	log := logger.New()

	// Настройка JWT секрета
	jwt.SetSecret(cfg.JWT.SecretKey)

	// Подключение к базе данных
	connstr := cfg.DB.GetConnStr()
	repo, err := repository.New(connstr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer repo.Close()

	// Применение миграций
	migration.ApplyMigrations(connstr)

	// Инициализация Redis кэша
	redis := cache.NewRedisCache(cfg.Redis.Host, cfg.Redis.Port)

	// Инициализация балансировщика
	bal := balancer.NewBalancer(cfg.Backends)

	// Запуск проверки состояния сервиса
	go health.HealthChecker(bal)

	// Инициализация лимитера
	limiter := limiter.New()

	// Создание сервисов
	svc := service.New(repo, redis, limiter)

	// Инициализация обработчиков
	h := handler.New(svc, bal, log)

	// Инициализация роутера и настройка маршрутов
	r := router.New()
	h.SetupRoutes(r)

	// Запуск HTTP сервера
	server.Start(server.HTTPServerConfig{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		Timeout:      cfg.Server.Timeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}, log)
}
