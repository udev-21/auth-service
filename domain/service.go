package domain

import (
	"context"
	"strings"
	"udev21/auth/config"
	"udev21/auth/domain/base"
	myErrors "udev21/auth/error"
	"udev21/auth/util"
)

type Service struct {
	base.Model
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	OwnerID     string  `json:"owner_id" db:"owner_id"`
}

func (s Service) GetMysqlTableName() string {
	return "services"
}

func (s *Service) ValidateName() error {
	s.Name = strings.TrimSpace(s.Name)
	if len(s.Name) < 1 {
		return myErrors.ErrInvalidValue
	}
	return nil
}

func (s *Service) ValidateDescription() error {
	if s.Description != nil {
		*s.Description = strings.TrimSpace(*s.Description)
		if len(*s.Description) < 1 {
			return myErrors.ErrInvalidValue
		}
	}
	return nil
}

func (s *Service) ValidateOwnerID() error {
	if len(s.OwnerID) < config.DBTableIDKeyLength {
		return myErrors.ErrInvalidValue
	}
	return nil
}

type IServiceUseCase interface {
	CreateService(ctx context.Context, owner *User, service *ServiceCreateInput) (*Service, error)
	GetServiceByID(ctx context.Context, id string) (*Service, error)
	GetServicesByOwner(ctx context.Context, owner *User) ([]Service, error)
	AddExistingUserToService(ctx context.Context, user *User, service *Service) error
}

type IServiceRepository interface {
	GetServiceByID(ctx context.Context, id string) (*Service, error)
	GetAllServiceByOwnerID(ctx context.Context, ownerID string) ([]Service, error)
	GetServiceByPosition(ctx context.Context, position uint64) (*Service, error)
	AddExistingUserToService(ctx context.Context, user *User, service *Service) error
	Create(ctx context.Context, service *ServiceCreateInput) (*Service, error)
}

type ServiceCreateInput struct {
	Service
}

func (s *ServiceCreateInput) Validate() error {
	return util.GetErrorIfExist(s.ValidateName, s.ValidateDescription)
}

type CreateServiceUserInput struct {
	UserCreateInput
	ServiceID string `json:"service_id"`
}

func (s *CreateServiceUserInput) ValidateServiceID() error {
	if len(s.ServiceID) < config.DBTableIDKeyLength {
		return myErrors.ErrInvalidValue
	}
	return nil
}

func (s *CreateServiceUserInput) Validate() error {
	return util.GetErrorIfExist(s.ValidateServiceID, s.UserCreateInput.Validate)
}
