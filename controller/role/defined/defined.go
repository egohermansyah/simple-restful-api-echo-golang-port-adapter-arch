package defined

import (
	"database/sql"
	"golang-vscode-setup/service/role/defined"
	"time"
)

type InsertRequest struct {
	Name string         `json:"name" validate:"required,min=1"`
	Desc sql.NullString `json:"desc"`
}

type DefaultResponse struct {
	Id       uint      `json:"id"`
	Name     string    `json:"name"`
	Desc     string    `json:"desc"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"updated"`
}

func NewDefaultResponse(role *defined.Role) *DefaultResponse {
	return &DefaultResponse{
		Id:       role.Id,
		Name:     role.Name,
		Desc:     role.Desc.String,
		Created:  role.Created,
		Modified: role.Modified,
	}
}
