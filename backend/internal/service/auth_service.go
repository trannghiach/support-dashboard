package service

import (
	"context"
	"errors"
	
	"github.com/trannghiach/support-dashboard/backend/internal/auth"
	"github.com/trannghiach/support-dashboard/backend/internal/dto"
	"github.com/trannghiach/support-dashboard/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(s.jwtSecret, user.ID, user.Role)
	if err != nil {
		return "", err
	}
	
	return token, nil
}