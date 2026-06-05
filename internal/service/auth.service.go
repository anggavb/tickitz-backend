package service

import (
	"context"
	"log"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/errs"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/pkg"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {

	isEmailExists, err := s.authRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[Register] FindByEmail error: %v\n", err)
		return errs.ErrInternalServer
	}

	if isEmailExists {
		log.Printf("[Register] Email already exists: %s\n", req.Email)
		return errs.ErrExistingEmail
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	hashedPassword := hc.Hash(req.Password)

	if _, err := s.authRepo.Create(ctx, req.Email, hashedPassword); err != nil {
		log.Printf("[Register] Create user error: %v\n", err)
		return errs.ErrInternalServer
	}

	return nil
}
