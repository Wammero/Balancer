package models

import "time"

type TokenBucket struct {
	Key             string    `json:"key"`               // IP
	Tokens          int       `json:"tokens"`            // Текущее количество токенов
	LastRefill      time.Time `json:"last_refill"`       // Когда последний раз добавляли токены
	TokensPerSecond int       `json:"tokens_per_second"` // Скорость пополнения
	Capacity        int       `json:"capacity"`          // Максимум токенов
	CreatedAt       time.Time `json:"created_at"`        // Когда добавлен
	UpdatedAt       time.Time `json:"updated_at"`        // Когда обновлён
}
