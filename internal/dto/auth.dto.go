package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	IsAgree  bool   `json:"is_agree" binding:"required,eq=true"`
}
<<<<<<< HEAD

type ActivationRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type NewOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}
=======
>>>>>>> b9ee6f3b7daa7e17199dec072791cf7dbe5d369b
