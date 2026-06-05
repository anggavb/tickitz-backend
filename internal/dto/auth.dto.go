package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	IsAgree  bool   `json:"is_agree" binding:"required,eq=true"`
}
