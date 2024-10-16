package models

type User struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	IsVerified          bool   `json:"is_verified"`
	FailedLoginAttempts int
	IsBlocked           bool
}
