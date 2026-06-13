package controller

import (
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/jwttoken"
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

// GetProfileById godoc
// @Summary Get user profile
// @Description Get profile of currently logged in user
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserProfile
// @Failure 500 {object} dto.ErrorResponse
// @Router /profile [get]
func (c *ProfileController) GetProfileById(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := claims.UserId

	profile, err := c.ProfileService.GetUserProfile(
		ctx.Request.Context(),
		userID,
	)
	if err != nil {
		log.Printf("[GetProfileById] GetUserProfile error: %v\n", err)

		response.Error(
			ctx,
			http.StatusInternalServerError,
			"failed to get profile",
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success to get profile",
		profile,
	)
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update profile of currently logged in user
// @Tags Profile
// @Accept multipart/form-data
// @Param first_name formData string false "First Name"
// @Param last_name formData string false "Last Name"
// @Param phone formData string false "Phone Number"
// @Param photo formData file false "Profile Photo"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /profile [put]
func (c *ProfileController) UpdateUserProfile(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID := claims.UserId

	var req dto.UpdateProfileRequest

	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("[UpdateUserProfile] Bind error: %v\n", err)

		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var photo *multipart.FileHeader

	file, err := ctx.FormFile("photo")
	if err == nil {
		photo = file
	}

	if err := c.ProfileService.ChangeUserProfile(
		ctx.Request.Context(),
		req,
		photo,
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
