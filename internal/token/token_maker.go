package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type TokenMaker interface {
	CreateToken(username string, userID int64, duration time.Duration) (string, error)
	VerifyToken(token string) (*TokenPayload, error)
}

func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiresAt.Time) {
		return fmt.Errorf("token has expired")
	}
	return nil
}
