package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID       string `json:"user_id"`
	PlatformRole string `json:"platform_role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService() (*TokenService, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	ttl := 24 * time.Hour
	if v := os.Getenv("JWT_TTL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("JWT_TTL: %w", err)
		}
		ttl = d
	}

	return &TokenService{secret: []byte(secret), ttl: ttl}, nil
}

func (s *TokenService) Issue(userID uuid.UUID, platformRole string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:       userID.String(),
		PlatformRole: platformRole,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *TokenService) Parse(accessToken string) (uuid.UUID, string, error) {
	parsed, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return uuid.Nil, "", err
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return uuid.Nil, "", fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, "", err
	}

	return userID, claims.PlatformRole, nil
}
