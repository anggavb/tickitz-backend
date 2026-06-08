package service

import (
	"context"
	"log"
	"time"

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
		log.Printf("[Register] Email not activated: %s\n", req.Email)
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

	otp, hashedOTP := s.generateOTP()

	if _, err := s.authRepo.Create(
		ctx,
		req.Email,
		hashedPassword,
		hashedOTP,
	); err != nil {
		log.Printf("[Register] Create user error: %v\n", err)
		return errs.ErrInternalServer
	}

	if err := pkg.SendMail(
		[]string{req.Email},
		"[TICKITZ] Activation OTP",
		"Ini adalah kode OTP anda : \n\n"+otp,
	); err != nil {
		log.Printf("[Register] Send email error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Printf("[Register] OTP sent successfully to %s\n", req.Email)

	return nil
}

func (s *AuthService) Activate(ctx context.Context, req dto.ActivationRequest) error {

	isAccountNotActive, err := s.authRepo.FindByEmailAndActivate(ctx, req.Email)
	if err != nil {
		log.Printf("[Activate] FindByEmailAndActivate error: %v\n", err)
		return errs.ErrInternalServer
	}

	if !isAccountNotActive {
		log.Printf("[Activate] Account already activated: %s\n", req.Email)
		return errs.ErrAccountActivated
	}

	expiryToken, err := s.authRepo.GetExpiryToken(ctx, req.Email)
	if err != nil {
		log.Printf("[Activate] GetExpiryToken error: %v\n", err)
		return errs.ErrInternalServer
	}

	if time.Now().After(expiryToken) {
		log.Printf("[Activate] Token expired for %s\n", req.Email)
		return errs.ErrTokenExpired
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	existingToken, err := s.authRepo.GetUserToken(ctx, req.Email)
	if err != nil {
		log.Printf("[Activate] GetUserToken error: %v\n", err)
		return errs.ErrInternalServer
	}

	if err := hc.Compare(req.OTP, existingToken); err != nil {
		log.Printf("[Activate] Invalid OTP for %s\n", req.Email)
		return errs.ErrInvalidOTP
	}

	if err := s.authRepo.Activate(ctx, req.Email); err != nil {
		log.Printf("[Activate] Activate error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Printf("[Activate] Account activated successfully: %s\n", req.Email)

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

	otp, hashedOTP := s.generateOTP()

	if err := s.authRepo.UpdateOTP(ctx, req.Email, hashedOTP); err != nil {
		log.Printf("[GetNewOTP] UpdateOTP error: %v\n", err)
		return errs.ErrInternalServer
	}

	if err := pkg.SendMail(
		[]string{req.Email},
		"[TICKITZ] New Activation OTP",
		"Ini adalah kode OTP baru anda : \n\n"+otp,
	); err != nil {
		log.Printf("[GetNewOTP] Send email error: %v\n", err)
		return errs.ErrInternalServer
	}

	log.Printf("[GetNewOTP] OTP sent successfully to %s\n", req.Email)

	return nil
}

func (s *AuthService) generateOTP() (string, string) {
	var hc pkg.HashConfig

	hc.OwaspRecomendedHashConfig()

	otp := pkg.GenerateOTP()
	hashedOTP := hc.Hash(otp)

	return otp, hashedOTP
}
