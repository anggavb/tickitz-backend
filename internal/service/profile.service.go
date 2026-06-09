package service

import (
	"context"
	"log"

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
		log.Printf("[Get Profile] GetProfile error: %v\n", err)
		return dto.UserProfile{}, err
	}
	userProfile := dto.UserProfile{
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		Phone:       profile.Phone,
		Photo:       profile.Photo,
		LoyaltyTier: profile.LoyaltyTier,
	}
	return userProfile, nil
}
