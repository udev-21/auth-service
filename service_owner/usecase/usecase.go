package usercase

import (
	"context"
	"udev21/auth/domain"
)

type serviceOwnerUseCase struct {
	repo domain.IServiceOwnerRepository
}

func New(repo domain.IServiceOwnerRepository) domain.IServiceOwnerUseCase {
	return &serviceOwnerUseCase{
		repo: repo,
	}
}

func (u *serviceOwnerUseCase) IsServiceOwner(ctx context.Context, user *domain.User) (bool, error) {
	if user == nil {
		return false, nil
	}
	_, err := u.repo.GetServiceOwnerByID(ctx, user.ID)
	if err != nil {
		return false, err
	}

	return err == nil, nil
}
