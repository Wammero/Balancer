package jwt

import (
	"context"
	"time"

	"github.com/Wammero/Balancer/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret []byte

func SetSecret(secret string) {
	jwtSecret = []byte(secret)
}

func GetSecret() []byte {
	return jwtSecret
}

func GenerateJWT(clientID string) (string, error) {
	claims := &models.Claims{
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			Issuer:    "balancer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(models.ClientIDContextKey).(string)
	return id, ok
}
