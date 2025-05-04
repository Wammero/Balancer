package models

type Client struct {
	ClientID   string `json:"client_id"`    // ID
	Capacity   int    `json:"capacity"`     // Максимум токенов
	RatePerSec int    `json:"rate_per_sec"` // Скорость пополнения
}
