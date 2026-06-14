package model

type UserProfile struct {
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Phone       string `db:"phone"`
	Photo       string `db:"photo"`
	Point       int    `db:"point"`
	NextPoint   int    `db:"next_point"`
	NextTier    string `db:"next_tier"`
	LoyaltyTier string `db:"loyalty_tier"`
	Email       string `db:"email"`
}
