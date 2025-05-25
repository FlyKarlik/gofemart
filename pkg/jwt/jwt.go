package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type JWTPayload struct {
	SecretKey      string
	Issuer         string
	AccessTokenTTL time.Duration
	Claims         Claims
}

type Claims struct {
	UserID string `json:"user_id"`
	Login  string `json:"login,omitempty"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(payload JWTPayload) (string, error) {
	payload.Claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(payload.AccessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    payload.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload.Claims)
	return token.SignedString([]byte(payload.SecretKey))
}

func ParseToken(tokenString string, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func GetClearToken(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", ErrInvalidToken
	}
	return authHeader[7:], nil
}

func IsTokenExpired(err error) bool {
	return errors.Is(err, ErrExpiredToken)
}
