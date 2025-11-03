package jwt

import (
	"fmt"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider interface {
	GenerateToken(userID int64, userRole string, ttl time.Duration) (string, error)

	ParseToken(tokenStr string) (int64, string, error)
}

type jwtProviderImpl struct {
	secret string
}

func NewJWTProvider(secret string) JWTProvider {
	return &jwtProviderImpl{secret}
}

func (j *jwtProviderImpl) GenerateToken(userID int64, userRole string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": userRole,
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *jwtProviderImpl) ParseToken(tokenStr string) (int64, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return 0, "", common.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", common.ErrInvalidToken
	}

	idFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, "", common.ErrInvalidToken
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", common.ErrInvalidToken
	}

	return int64(idFloat), role, nil
}
