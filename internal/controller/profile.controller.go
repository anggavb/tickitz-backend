package controller

import (
	"github.com/gin-gonic/gin"
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
	// id ngambil dari claims -> menunggu middleware
}
