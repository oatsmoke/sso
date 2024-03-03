package user_profile_service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
	"sso/internal/model"
)

type UserProfileService struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewUserProfileService(log *slog.Logger, db *gorm.DB) *UserProfileService {
	return &UserProfileService{log: log, db: db}
}

func (s *UserProfileService) Create(ctx context.Context, login, password, email string) (int64, error) {
	const fn = "user_profile.service.Create"
	log := s.log.With(slog.String("fn", fn))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate password hash", slog.String("error:", err.Error()))
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	user := model.User{Login: login, PasswordHash: string(passHash), Email: email}

	if result := s.db.Create(&user); result.Error != nil {
		log.Error("failed to create user", slog.String("error:", result.Error.Error()))
		return 0, nil
	}

	log.Info("user registered")

	return user.ID, nil
}

func (s *UserProfileService) Authentication(ctx context.Context, login, password string) int64 {
	const fn = "user_profile.service.Authentication"
	log := s.log.With(slog.String("fn", fn))

	user := model.User{}

	s.db.Find(&user, "login = ?", login)

	if user.ID == 0 {
		log.Warn("user is not found")
		return 0
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Warn("wrong password", slog.String("error:", err.Error()))
		return 0
	}

	log.Info("user is authenticated")

	return user.ID
}

func (s *UserProfileService) Info(ctx context.Context, userId int64) (*model.User, error) {
	const fn = "user_profile.service.Info"
	log := s.log.With(slog.String("fn", fn))

	var user model.User

	tx := s.db.Find(&user, "id = ?", userId)

	if tx.Error != nil {
		log.Error("query find user", slog.String("error", tx.Error.Error()))
		return nil, tx.Error
	}

	if user.ID == 0 {
		log.Warn("user is not found")
		return nil, nil
	}

	log.Info("user found")

	return &user, nil
}
