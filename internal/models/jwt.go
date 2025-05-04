package models

import "github.com/golang-jwt/jwt/v4"

type ContextKey string

const (
	ClientIDContextKey ContextKey = "client_id"
)

type Claims struct {
	ClientID string `json:"client_id"`
	jwt.RegisteredClaims
}
