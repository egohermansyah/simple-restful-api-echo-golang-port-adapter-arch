package user

import (
	"net/url"
	"simple-restful-api-echo-golang-port-adapter-archcore/domains/user"
	port "simple-restful-api-echo-golang-port-adapter-archcore/ports/user"
	"simple-restful-api-echo-golang-port-adapter-archcore/services/util/generatepassword"
	"simple-restful-api-echo-golang-port-adapter-archutil/logger"
)

const LOG_IDENTIFIER = "SERVICE_USER"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

type Service struct {
	repository port.Repository
}

func New(repository port.Repository) port.Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) List(filters url.Values, limit int, offset int) ([]*user.User, error) {
	result, err := service.repository.List(filters, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) Create(data user.User) (*user.User, error) {
	password, err := generatepassword.GeneratePassword(data.Password, data.PasswordSalt)
	if err != nil {
		return nil, err
	}
	data.Password = password
	result, err := service.repository.Create(data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) FindById(id string) (*user.User, error) {
	result, err := service.repository.FindById(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) UpdateById(id string, data user.User) (*user.User, error) {
	result, err := service.repository.UpdateById(id, data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
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
