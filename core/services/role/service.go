package role

import (
	"net/url"
	"simple-restful-api-echo-golang-port-adapter-archcore/domains/role"
	port "simple-restful-api-echo-golang-port-adapter-archcore/ports/role"
	"simple-restful-api-echo-golang-port-adapter-archutil/logger"
)

const LOG_IDENTIFIER = "SERVICE_ROLE"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

type Service struct {
	repository port.Repository
}

func New(repository port.Repository) port.Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) List(filters url.Values, limit int, offset int) ([]*role.Role, error) {
	result, err := service.repository.List(filters, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) FindById(id string) (*role.Role, error) {
	result, err := service.repository.FindById(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) UpdateById(id string, data role.Role) (*role.Role, error) {
	result, err := service.repository.UpdateById(id, data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}
