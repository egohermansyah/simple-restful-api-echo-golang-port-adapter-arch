package role

import (
	"net/url"
	"simple-restful-api-echo-golang-port-adapter-archcore/domains/role"
)

type Repository interface {
	List(filters url.Values, limit int, offset int) ([]*role.Role, error)
	Create(data role.Role) (*role.Role, error)
	FindById(id string) (*role.Role, error)
	UpdateById(id string, data role.Role) (*role.Role, error)
	DeleteById(id string) error
}

type Service interface {
	List(filters url.Values, limit int, offset int) ([]*role.Role, error)
	FindById(id string) (*role.Role, error)
	UpdateById(id string, data role.Role) (*role.Role, error)
}
