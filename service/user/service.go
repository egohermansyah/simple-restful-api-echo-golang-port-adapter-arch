package user

import (
	"golang-vscode-setup/repository/user"
	"golang-vscode-setup/service/user/defined"
	"golang-vscode-setup/service/util/generatepassword"
	"golang-vscode-setup/util/logger"
	"net/url"
)

const LOG_IDENTIFIER = "SERVICE_USER"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func NewUser(user defined.User) *defined.User {
	return &defined.User{
		Id:           user.Id,
		RoleId:       user.RoleId,
		Email:        user.Email,
		Password:     user.Password,
		PasswordSalt: user.PasswordSalt,
		PhoneNumber:  user.PhoneNumber,
		LoginAttempt: user.LoginAttempt,
		IsLogin:      user.IsLogin,
		IsActive:     user.IsActive,
		Created:      user.Created,
		Modified:     user.Modified,
	}
}

type Service struct {
	repository user.IRepository
}

type IService interface {
	List(filters url.Values, limit int, offset int) ([]*defined.User, error)
	Create(role defined.User) (*defined.User, error)
	FindById(id string) (*defined.User, error)
	UpdateById(id string, role defined.User) (*defined.User, error)
	DeleteById(id string) error
}

func NewService(repository user.IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (service *Service) List(filters url.Values, limit int, offset int) ([]*defined.User, error) {
	result, err := service.repository.List(filters, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) Create(user defined.User) (*defined.User, error) {
	password, err := generatepassword.GeneratePassword(user.Password, user.PasswordSalt)
	if err != nil {
		return nil, err
	}
	user.Password = password
	repository, err := service.repository.Create(user)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewUser(*repository)
	return result, nil
}

func (service *Service) FindById(id string) (*defined.User, error) {
	repository, err := service.repository.FindById(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewUser(*repository)
	return result, nil
}

func (service *Service) UpdateById(id string, user defined.User) (*defined.User, error) {
	repository, err := service.repository.UpdateById(id, user)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewUser(*repository)
	return result, nil
}

func (service *Service) DeleteById(id string) error {
	err := service.repository.DeleteById(id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
