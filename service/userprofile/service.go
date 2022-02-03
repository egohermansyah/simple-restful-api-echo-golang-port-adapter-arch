package userprofile

import (
	"golang-vscode-setup/repository/userprofile"
	"golang-vscode-setup/service/userprofile/defined"
	"golang-vscode-setup/util/logger"
)

const LOG_IDENTIFIER = "SERVICE_CUSTOMER"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func NewUserProfile(userProfile defined.UserProfile) *defined.UserProfile {
	return &defined.UserProfile{
		UserId:     userProfile.UserId,
		FirstName:  userProfile.FirstName,
		LastName:   userProfile.LastName,
		Otp:        userProfile.Otp,
		OtpCreated: userProfile.OtpCreated,
		FcmToken:   userProfile.FcmToken,
	}
}

type Service struct {
	repository userprofile.IRepository
}

type IService interface {
	Create(userProfile defined.UserProfile) (*defined.UserProfile, error)
	FindById(userId string) (*defined.UserProfile, error)
	UpdateById(userId string, userProfile defined.UserProfile) (*defined.UserProfile, error)
	DeleteById(userId string) error
}

func NewService(repository userprofile.IRepository) IService {
	return &Service{
		repository: repository,
	}
}

func (service *Service) Create(userProfike defined.UserProfile) (*defined.UserProfile, error) {
	result, err := service.repository.Create(userProfike)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) FindById(userId string) (*defined.UserProfile, error) {
	result, err := service.repository.FindById(userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) UpdateById(userId string, userProfile defined.UserProfile) (*defined.UserProfile, error) {
	result, err := service.repository.UpdateById(userId, userProfile)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, nil
}

func (service *Service) DeleteById(userId string) error {
	err := service.repository.DeleteById(userId)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
