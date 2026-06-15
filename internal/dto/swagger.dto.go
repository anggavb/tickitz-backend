package dto

type EmptyDataResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data"`
}

type AuthEmailData struct {
	Email string `json:"email" example:"user@example.com"`
}

type AuthEmailSuccessResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Registration successful. Please check your email to activate your account."`
	Data    AuthEmailData `json:"data"`
}

type LoginSuccessResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Login success"`
	Data    LoginResponse `json:"data"`
}

type StringListSuccessResponse struct {
	Success bool     `json:"success" example:"true"`
	Message string   `json:"message" example:"success"`
	Data    []string `json:"data"`
}

type CinemaListSuccessResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"success"`
	Data    []CinemaResponse `json:"data"`
}

type MovieScheduleListSuccessResponse struct {
	Success bool                    `json:"success" example:"true"`
	Message string                  `json:"message" example:"movie schedules retrieved successfully"`
	Data    []MovieScheduleResponse `json:"data"`
}

type MovieLocationListSuccessResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"movie locations retrieved successfully"`
	Data    []MovieLocationRow `json:"data"`
}

type MovieShowtimeListSuccessResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"movie showtimes retrieved successfully"`
	Data    []MovieShowtimeRow `json:"data"`
}

type MovieScheduleOptionsSuccessResponse struct {
	Success bool                         `json:"success" example:"true"`
	Message string                       `json:"message" example:"movie schedule options retrieved successfully"`
	Data    MovieScheduleOptionsResponse `json:"data"`
}

type MovieDetailSuccessResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"movie detail retrieved successfully"`
	Data    MovieDetailResponse `json:"data"`
}

type MovieHomeListSuccessResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:"success to get movies"`
	Data    GetAllMoviesResponse `json:"data"`
}

type UpcomingMoviesSuccessResponse struct {
	Success bool                   `json:"success" example:"true"`
	Message string                 `json:"message" example:"success to get upcoming movies"`
	Data    []MoviePreviewResponse `json:"data"`
}

type CreatePendingOrderSuccessResponse struct {
	Success bool                       `json:"success" example:"true"`
	Message string                     `json:"message" example:"order created successfully"`
	Data    CreatePendingOrderResponse `json:"data"`
}

type OrderDetailSuccessResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"order detail retrieved successfully"`
	Data    OrderDetailResponse `json:"data"`
}

type PaymentMethodsSuccessResponse struct {
	Success bool                         `json:"success" example:"true"`
	Message string                       `json:"message" example:"payment methods retrieved successfully"`
	Data    []OrderPaymentMethodResponse `json:"data"`
}

type OrderHistorySuccessResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:"order history retrieved successfully"`
	Data    OrderHistoryResponse `json:"data"`
}

type UserProfileSuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"success to get profile"`
	Data    UserProfile `json:"data"`
}
