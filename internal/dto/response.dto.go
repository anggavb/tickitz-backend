package dto

type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"date format must be YYYY-MM-DD"`
}
