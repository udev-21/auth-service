package domain

import (
	"context"
	"udev21/auth/domain/base"
)

type ServiceOwner struct {
	base.Model
	UserId string `json:"user_id" db:"user_id"`
}

func (s ServiceOwner) GetMysqlTableName() string {
	return "service_owners"
}

type IServiceOwnerUseCase interface {
	IsServiceOwner(ctx context.Context, user *User) (bool, error)
	HasService(ctx context.Context, user *User, service *Service) (bool, error)
	// HasUser(ctx context.Context, ownerUser, user *User) (bool, error)
}

type IServiceOwnerRepository interface {
	GetServiceOwnerByID(ctx context.Context, id string) (*ServiceOwner, error)
	HasService(ctx context.Context, user *User, service *Service) (bool, error)
	// HasUser(ctx context.Context, ownerUser, user *User) (bool, error)

}
