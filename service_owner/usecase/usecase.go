package usercase

import (
	"context"
	"udev21/auth/domain"
	myErrors "udev21/auth/error"
)

type serviceOwnerUseCase struct {
	serviceOwnerRepo domain.IServiceOwnerRepository
	serviceRepo      domain.IServiceRepository
}

func New(serviceOwnerRepo domain.IServiceOwnerRepository, serviceRepo domain.IServiceRepository) domain.IServiceOwnerUseCase {
	return &serviceOwnerUseCase{
		serviceOwnerRepo: serviceOwnerRepo,
		serviceRepo:      serviceRepo,
	}
}

func (u *serviceOwnerUseCase) IsServiceOwner(ctx context.Context, user *domain.User) (bool, error) {
	if user == nil {
		return false, nil
	}
	_, err := u.serviceOwnerRepo.GetServiceOwnerByID(ctx, user.ID)
	if err != nil {
		return false, err
	}

	return err == nil, nil
}

func (u *serviceOwnerUseCase) CreateService(ctx context.Context, owner *domain.User, service *domain.ServiceCreateInput) (*domain.Service, error) {
	if service.Validate() != nil {
		return nil, myErrors.ErrInvalidInput
	}
	ok, err := u.IsServiceOwner(ctx, owner)
	if err != nil || !ok {
		return nil, myErrors.ErrNotFound
	}
	service.OwnerID = owner.ID

	return u.serviceRepo.Create(ctx, service)
}

// GetServices
func (u *serviceOwnerUseCase) GetServices(ctx context.Context, owner *domain.User) ([]domain.Service, error) {
	ok, err := u.IsServiceOwner(ctx, owner)
	if err != nil || !ok {
		return nil, myErrors.ErrNotFound
	}
	return u.serviceRepo.GetAllServiceByOwnerID(ctx, owner.ID)
}

// func (r *serviceOwnerUseCase) HasUser(ctx context.Context, ownerUser, user *domain.User) (bool, error) {
// 	// return r.serviceOwnerRepo.HasUser(ctx, ownerUser, user)
// }

// func (r *serviceOwnerUseCase) AddExistingUserToService(ctx context.Context, user *domain.User, service *domain.Service) error {

// 	return nil
// }
