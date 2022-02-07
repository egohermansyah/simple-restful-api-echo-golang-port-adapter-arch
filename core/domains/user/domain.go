package user

import (
	"time"
)

type User struct {
	Id           uint      `json:"id"`
	RoleId       uint      `json:"role_id"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	PasswordSalt string    `json:"password_salt"`
	PhoneNumber  string    `json:"phone_number"`
	LoginAttempt uint8     `json:"login_attempt"`
	IsLogin      bool      `json:"is_login"`
	IsActive     bool      `json:"is_active"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
}
