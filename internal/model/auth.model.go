package model

type GetUserLogin struct {
	Id         int    `db:"id"`
	Password   string `db:"password"`
	Photo      string `db:"photo"`
	Role       string `db:"role"`
	IsVerified bool   `db:"is_verified"`
}
