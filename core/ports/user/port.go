package user

import (
	"net/url"
	"simple-restful-api-echo-golang-port-adapter-archcore/domains/user"
)

type Repository interface {
	List(filters url.Values, limit int, offset int) ([]*user.User, error)
	Create(data user.User) (*user.User, error)
	FindById(id string) (*user.User, error)
	UpdateById(id string, data user.User) (*user.User, error)
	DeleteById(id string) error
}

type Service interface {
	List(filters url.Values, limit int, offset int) ([]*user.User, error)
	Create(data user.User) (*user.User, error)
	FindById(id string) (*user.User, error)
	UpdateById(id string, data user.User) (*user.User, error)
	DeleteById(id string) error
}
