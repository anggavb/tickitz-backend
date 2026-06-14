package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/repository"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

func (s *ProfileService) GetUserProfile(ctx context.Context, id int) (dto.UserProfile, error) {
	profile, err := s.profileRepo.GetProfile(ctx, id)
	if err != nil {
		log.Printf("[GetUserProfile] GetProfile error: %v\n", err)
		return dto.UserProfile{}, err
	}

	return dto.UserProfile{
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		Phone:       profile.Phone,
		Photo:       profile.Photo,
		Point:       profile.Point,
		NextTier:    profile.NextTier,
		NextPoint:   profile.NextPoint,
		LoyaltyTier: profile.LoyaltyTier,
		Email:       profile.Email,
	}, nil
}

func (s *ProfileService) ChangeUserProfile(
	ctx context.Context,
	req dto.UpdateProfileRequest,
	photo *multipart.FileHeader,
	userID int,
) error {

	if photo != nil {
		ext := filepath.Ext(photo.Filename)

		filename := fmt.Sprintf(
			"profile_%d_%d%s",
			userID,
			time.Now().Unix(),
			ext,
		)

		dst := filepath.Join("public/img/profile", filename)

		src, err := photo.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, src); err != nil {
			return err
		}

		photoURL := "/profile/" + filename

		req.Photo = &photoURL
	}

	if err := s.profileRepo.UpdateProfile(
		ctx,
		req,
		userID,
	); err != nil {
		log.Printf("[ChangeUserProfile] UpdateProfile error: %v\n", err)
		return err
	}

	return nil
}
