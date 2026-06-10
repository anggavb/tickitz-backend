package dto

type UserProfile struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Photo       string `json:"photo"`
	LoyaltyTier string `json:"loyalty_tier"`
}

type UpdateProfileRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	Photo     *string `json:"photo"`
}
