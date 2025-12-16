package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtMaker struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("secret key must be at least 32 characters")
	}
	return &JwtMaker{secretKey: secretKey}, nil
}
func (maker *JwtMaker) CreateToken(username string, userID int64, duration time.Duration) (string, error) {
	payload := &TokenPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			Issuer:    "gobank",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserID:   userID,
		Username: username,
	}
	tokenjwt := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := tokenjwt.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
func (maker *JwtMaker) VerifyToken(token string) (*TokenPayload, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &TokenPayload{}, keyfunc)
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	payload, ok := jwtToken.Claims.(*TokenPayload)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return payload, nil
}
