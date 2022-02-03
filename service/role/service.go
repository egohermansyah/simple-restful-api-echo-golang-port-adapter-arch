package role

import (
	"golang-vscode-setup/repository/role"
	"golang-vscode-setup/service/role/defined"
	"golang-vscode-setup/util/logger"
	"net/url"
)

const LOG_IDENTIFIER = "SERVICE_ROLE"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func NewRole(role defined.Role) *defined.Role {
	return &defined.Role{
		Id:       role.Id,
		Name:     role.Name,
		Desc:     role.Desc,
		Created:  role.Created,
		Modified: role.Modified,
	}
}

type Service struct {
	repository role.IRepository
}

type IService interface {
	List(filters url.Values, limit int, offset int) ([]*defined.Role, error)
	FindById(id string) (*defined.Role, error)
	UpdateById(id string, role defined.Role) (*defined.Role, error)
}

func NewService(repository role.IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (service *Service) List(filters url.Values, limit int, offset int) ([]*defined.Role, error) {
	result, err := service.repository.List(filters, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) FindById(id string) (*defined.Role, error) {
	repository, err := service.repository.FindById(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewRole(*repository)
	return result, nil
}

func (service *Service) UpdateById(id string, role defined.Role) (*defined.Role, error) {
	repository, err := service.repository.UpdateById(id, role)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := NewRole(*repository)
	return result, nil
}
