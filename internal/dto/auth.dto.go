package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	IsAgree  bool   `json:"is_agree" binding:"required,eq=true"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type ActivationRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type NewOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginResponse struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
	Photo string `json:"photo"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
