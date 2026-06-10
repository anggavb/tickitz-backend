package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/response"
	"github.com/tickitz-backend/internal/service"
)

type ProfileController struct {
	ProfileService *service.ProfileService
}

func NewProfileController(profileService *service.ProfileService) *ProfileController {
	return &ProfileController{
		ProfileService: profileService,
	}
}

func (c *ProfileController) GetProfileById(ctx *gin.Context) {
	userID := 2

	profile, err := c.ProfileService.GetUserProfile(
		ctx.Request.Context(),
		userID,
	)
	if err != nil {
		log.Printf("[GetProfileById] GetUserProfile error: %v\n", err)
		response.Error(ctx, http.StatusInternalServerError, "failed to get profile")
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success to get profile",
		profile,
	)
}

func (c *ProfileController) UpdateUserProfile(ctx *gin.Context) {
	// TODO: ambil dari JWT claims
	userID := 2

	var req dto.UpdateProfileRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateUserProfile] BindJSON error: %v\n", err)

		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.ProfileService.ChangeUserProfile(
		ctx.Request.Context(),
		req,
		userID,
	); err != nil {
		log.Printf("[UpdateUserProfile] Service error: %v\n", err)

		response.Error(
			ctx,
			http.StatusInternalServerError,
			"failed to update profile",
		)
		return
	}

	response.Success(ctx, http.StatusOK, "profile updated", nil)
}
