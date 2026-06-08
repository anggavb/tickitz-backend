package service

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
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
		token := uuid.NewString()

		if _, err := s.authRepo.Create(ctx, req.Email, hashedPassword, token); err != nil {
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
	return nil

}

func (s *AuthService) Activate(ctx context.Context, req dto.ActivationRequest) error {

	isAccountNotActive, err := s.authRepo.FindByEmailAndActivate(ctx, req.Email)
	if err != nil {
		log.Printf("[Activate] FindByEmailAndActivate error: %v", err)
		return errs.ErrInternalServer
	}

	if !isAccountNotActive {
		log.Printf("[Activate] Account already activated: %s", req.Email)
		return errs.ErrAccountActivated
	}

	expiryToken, err := s.authRepo.GetExpiryToken(ctx, req.Email)
	if err != nil {
		log.Printf("[Activate] GetExpiryToken error: %v", err)
		return errs.ErrInternalServer
	}

	if time.Now().After(expiryToken) {
		log.Printf("[Activate] Token expired for %s", req.Email)
		return errs.ErrTokenExpired
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	existingToken, err := s.authRepo.GetUserToken(ctx, req.Email)
	if err != nil {
		log.Printf("[Activate] GetUserToken error: %v", err)
		return errs.ErrInternalServer
	}

	if err := hc.Compare(req.OTP, existingToken); err != nil {
		log.Printf("[Activate] Invalid OTP for %s", req.Email)
		return errs.ErrInvalidOTP
	}

	if err := s.authRepo.Activate(ctx, req.Email); err != nil {
		log.Printf("[Activate] Activate error: %v", errs.ErrInvalidOTP)
		return errs.ErrInternalServer
	}

	log.Printf("[Activate] Account activated successfully: %s", req.Email)

	return nil
}

func (s *AuthService) GetNewOTP(ctx context.Context, req dto.NewOTPRequest) error {

	isAccountNotActive, err := s.authRepo.FindByEmailAndActivate(ctx, req.Email)
	if err != nil {
		log.Printf("[GetNewOTP] FindByEmailAndActivate error: %v\n", err)
		return errs.ErrInternalServer
	}

	if !isAccountNotActive {
		log.Printf("[GetNewOTP] Account already activated: %s\n", req.Email)
		return errs.ErrAccountActivated
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	otp := pkg.GenerateOTP()
	hashedOTP := hc.Hash(otp)

	if err := s.authRepo.UpdateOTP(ctx, req.Email, hashedOTP); err != nil {
		log.Printf("[GetNewOTP] UpdateOTP error: %v\n", err)
		return errs.ErrInternalServer
	}

	subject := "[TICKITZ] New Activation OTP"
	body := "Ini adalah kode OTP baru anda : \n\n" + otp

	if err := pkg.SendMail([]string{req.Email}, subject, body); err != nil {
		log.Printf("[GetNewOTP] Send email error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Printf("[GetNewOTP] New OTP sent successfully to %s\n", req.Email)

	return nil
}
