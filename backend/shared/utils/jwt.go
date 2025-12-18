package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidTokenExpired = errors.New("token expired")
	ErrInvalidTokenSigningMethod = errors.New("invalid token signing method")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
)
type JwtClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status"`
	jwt.RegisteredClaims
}

type JwtManager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJwtManager(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) *JwtManager {
	return &JwtManager{
		secretKey:            secretKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (m *JwtManager) AccessTokenDuration() time.Duration {
	return m.accessTokenDuration
}

func (m *JwtManager) RefreshTokenDuration() time.Duration {
	return m.refreshTokenDuration
}

func (m *JwtManager) GenerateAccessToken(userId, email, role, status string) (string, error) {
	claims := JwtClaims{
		UserID: userId,
		Email:  email,
		Role:   role,
		Status: status,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return tokenString, nil
}

func (m *JwtManager) GenerateRefreshToken(userId, email, role, status string) (string, error) {
	claims := JwtClaims{
		UserID: userId,
		Email:  email,
		Role:   role,
		Status: status,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return tokenString, nil
}

func (m *JwtManager) VerifyToken(tokenString string) (JwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, 
		&JwtClaims{},
		 func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidTokenSigningMethod
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return JwtClaims{}, ErrInvalidToken
	}
	claims, ok := token.Claims.(*JwtClaims)
	if !ok {
		return JwtClaims{}, ErrInvalidTokenClaims
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return JwtClaims{}, ErrInvalidTokenExpired
	}
	return *claims, nil
} 

