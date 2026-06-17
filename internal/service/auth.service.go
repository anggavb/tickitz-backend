package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/errs"
	"github.com/tickitz-backend/internal/jwttoken"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/pkg"
)

type AuthService struct {
	authRepo  *repository.AuthRepository
	authCache *repository.AuthCacheRepository
}

func NewAuthService(authRepo *repository.AuthRepository, authCache *repository.AuthCacheRepository) *AuthService {
	return &AuthService{
		authRepo:  authRepo,
		authCache: authCache,
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

	htmlBody := pkg.GenerateOTPEmail(otp)

	if err := pkg.SendMail(
		[]string{req.Email},
		"[TICKITZ] Activation OTP",
		htmlBody,
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

	htmlBody := pkg.GenerateResendOTPEmail(otp)

	if err := pkg.SendMail(
		[]string{req.Email},
		"[TICKITZ] New Activation OTP",
		htmlBody,
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

func (s *AuthService) Login(ctx context.Context, email, password string) (dto.LoginResponse, error) {
	user, err := s.authRepo.GetUserPassword(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.LoginResponse{}, errs.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, err
	}

	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	if err := hc.Compare(password, user.Password); err != nil {
		return dto.LoginResponse{}, errs.ErrInvalidCredentials
	}

	if !user.IsVerified {
		return dto.LoginResponse{}, errs.ErrAccountNotActive
	}

	claims := pkg.NewClaims(user.Id, email, user.Role)

	token, err := claims.GenJWT()
	if err != nil {
		log.Printf("[Login] Generate JWT error: %v\n", err)
		return dto.LoginResponse{}, errs.ErrInternalServer
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil || expiresAt == nil {
		log.Printf("[Login] Get JWT expiration error: %v\n", err)
		return dto.LoginResponse{}, errs.ErrInternalServer
	}

	if err := s.authCache.StoreToken(ctx, jwttoken.HashToken(token), user.Id, expiresAt.Time); err != nil {
		log.Printf("[Login] Store token error: %v\n", err)
		return dto.LoginResponse{}, errs.ErrInternalServer
	}

	data := dto.LoginResponse{
		Id:    user.Id,
		Photo: user.Photo,
		Token: token,
		Role:  user.Role,
	}

	return data, nil

}

func (s *AuthService) Logout(ctx context.Context, tokenHash string, userID int) error {
	if err := s.authCache.DeleteToken(ctx, tokenHash, userID); err != nil {
		log.Printf("[Logout] Delete token error: %v\n", err)
		return errs.ErrInternalServer
	}

	return nil
}

func (s *AuthService) ChangeUserPassword(ctx context.Context, newPassword string, id int) error {
	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()

	hashedPass := hc.Hash(newPassword)

	if err := s.authRepo.UpdatePassword(ctx, hashedPass, id); err != nil {
		log.Printf("[ChangeUserPassword] UpdatePassword error: %v\n", err)
		return err
	}

	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {

	userID, err := s.authRepo.GetUserIDByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	clientURL := os.Getenv("CLIENT_URL")

	rawToken := pkg.GenerateRandomToken(32)

	expiresAt := time.Now().Add(30 * time.Minute)

	// simpan ke redis
	err = s.authCache.StoreTokenForgotPassword(
		ctx,
		rawToken,
		int(userID),
		expiresAt,
	)

	if err != nil {
		return err
	}

	resetURL := fmt.Sprintf(
		"%s/auth/reset-password?token=%s",
		clientURL,
		rawToken,
	)

	htmlBody := pkg.GenerateForgotPasswordEmail(resetURL)

	err = pkg.SendMail(
		[]string{email},
		"[TICKITZ] Reset Password",
		htmlBody,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, password string) error {
	userRedisKey := fmt.Sprintf("tickitz:auth:reset-password:%s", token)
	exists, err := s.authCache.IsFogotPasswordKeyActive(ctx, userRedisKey)
	if err != nil {
		log.Printf("[ResetPassword] IsFogotPasswordKeyActive error: %v\n", err)
		return err
	}
	if !exists {
		log.Printf("[ResetPassword] IsExistToken error: %v\n%s", err, userRedisKey)
		return errors.New("expired token")
	}
	value, err := s.authCache.GetValueAndDelete(ctx, userRedisKey)
	if err != nil {
		log.Printf("[ResetPassword] GetValueAndDelete error: %v\n", err)
		return err
	}
	var hc pkg.HashConfig
	hc.OwaspRecomendedHashConfig()
	hashedNewPassword := hc.Hash(password)

	userID, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("[ResetPassword] strconv atoi error: %v\n", err)
		return err
	}

	if err := s.authRepo.UpdatePassword(ctx, hashedNewPassword, userID); err != nil {
		log.Printf("[ResetPassword] UpdatePassword error: %v\n", err)
		return err
	}
	return nil
}
