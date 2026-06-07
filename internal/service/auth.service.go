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

	isAccountNotActive, err := s.authRepo.FindByEmailAndActivate(ctx, req.Email)
	if err != nil {
		log.Printf("[Register] FindByEmailAndActivate error: %v\n", err)
		return errs.ErrInternalServer
	}

	if isAccountNotActive {
		log.Printf("[Register] Email haven't activate yet: %s\n", req.Email)
		return errs.ErrAccountNotActive
	}

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

	OTP := pkg.GenerateOTP()
	hashedOTP := hc.Hash(OTP)

	if _, err := s.authRepo.Create(ctx, req.Email, hashedPassword, hashedOTP); err != nil {
		log.Printf("[Register] Create user error: %v\n", err)
		return errs.ErrInternalServer
	}

	subject := "[TICKITZ] Activation Account"
	body := "Ini adalah kode OTP anda : \n\n" + OTP

	if err := pkg.SendMail([]string{req.Email}, subject, body); err != nil {
		log.Printf("[Register] Send email error: %v\n", err)
		return errs.ErrInternalServer
	}

	return nil
}
