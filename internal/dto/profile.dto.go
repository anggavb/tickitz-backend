package dto

type UserProfile struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Photo       string `json:"photo"`
	Point       int    `json:"point"`
	NextTier    string `json:"next_tier"`
	NextPoint   int    `json:"next_point"`
	LoyaltyTier string `json:"loyalty_tier"`
	Email       string `json:"email"`
}

type UpdateProfileRequest struct {
	FirstName *string `form:"first_name"`
	LastName  *string `form:"last_name"`
	Phone     *string `form:"phone"`

	Photo *string `form:"-"`
}
