package defined

import (
	"golang-vscode-setup/service/user/defined"
	"time"
)

type InsertRequest struct {
	RoleId      uint   `json:"role_id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateByIdRequest struct {
	RoleId      uint   `json:"role_id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number"`
	IsLogin     bool   `json:"is_login"`
	IsActive    bool   `json:"is_active"`
}

type DefaultResponse struct {
	Id           uint      `json:"id"`
	RoleId       uint      `json:"role_id"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	LoginAttempt uint8     `json:"login_attempt"`
	IsLogin      bool      `json:"is_login"`
	IsActive     bool      `json:"is_active"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
}

func NewDefaultResponse(user *defined.User) *DefaultResponse {
	return &DefaultResponse{
		Id:           user.Id,
		RoleId:       user.RoleId,
		Email:        user.Email.String,
		PhoneNumber:  user.PhoneNumber.String,
		LoginAttempt: user.LoginAttempt,
		IsLogin:      user.IsLogin,
		IsActive:     user.IsActive,
		Created:      user.Created,
		Modified:     user.Modified,
	}
}
