package auth_service

import (
	"context"
	"fmt"
	"log/slog"
	"sso/internal/lib/jwt"
	"time"
)

type AuthService struct {
	log      *slog.Logger
	tokenTTL time.Duration
	secret   string
}

func NewAuthService(log *slog.Logger, tokenTTL time.Duration, secret string) *AuthService {
	return &AuthService{
		log:      log,
		tokenTTL: tokenTTL,
		secret:   secret,
	}
}

func (s *AuthService) GetToken(ctx context.Context, userId int64) (string, error) {
	const fn = "auth_service_GetToken"
	log := s.log.With(slog.String("fn", fn))

	token, err := jwt.NewToken(userId, s.tokenTTL, s.secret)

	if err != nil {
		log.Error("failed to generate token", slog.String("error:", err.Error()))
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("token generated")

	return token, nil
}
