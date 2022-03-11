package usecase

import (
	"context"
	"udev21/auth/domain"
	myErrors "udev21/auth/error"
)

type serviceUseCase struct {
	serviceRepo domain.IServiceRepository
	userUseCase domain.IUserUseCase
}

func New(serviceRepo domain.IServiceRepository, userUseCase domain.IUserUseCase) domain.IServiceUseCase {
	return &serviceUseCase{
		serviceRepo: serviceRepo,
		userUseCase: userUseCase,
	}
}

func (u *serviceUseCase) GetServiceByID(ctx context.Context, id string) (*domain.Service, error) {
	return u.serviceRepo.GetServiceByID(ctx, id)
}

func (u *serviceUseCase) CreateService(ctx context.Context, owner *domain.User, service *domain.ServiceCreateInput) (*domain.Service, error) {
	if service.Validate() != nil {
		return nil, myErrors.ErrInvalidInput
	}
	service.OwnerID = owner.ID
	return u.serviceRepo.Create(ctx, service)
}

func (u *serviceUseCase) GetServicesByOwner(ctx context.Context, owner *domain.User) ([]domain.Service, error) {
	return u.serviceRepo.GetAllServiceByOwnerID(ctx, owner.ID)
}

func (u *serviceUseCase) AddExistingUserToService(ctx context.Context, user *domain.User, service *domain.Service) error {
	return u.serviceRepo.AddExistingUserToService(ctx, user, service)
}
