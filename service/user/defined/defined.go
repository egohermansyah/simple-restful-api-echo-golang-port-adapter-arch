package defined

import (
	"database/sql"
	"time"
)

type User struct {
	Id           uint           `json:"id"`
	RoleId       uint           `json:"role_id"`
	Email        sql.NullString `json:"email"`
	Password     string         `json:"password"`
	PasswordSalt string         `json:"password_salt"`
	PhoneNumber  sql.NullString `json:"phone_number"`
	LoginAttempt uint8          `json:"login_attempt"`
	IsLogin      bool           `json:"is_login"`
	IsActive     bool           `json:"is_active"`
	Created      time.Time      `json:"created"`
	Modified     time.Time      `json:"modified"`
}
