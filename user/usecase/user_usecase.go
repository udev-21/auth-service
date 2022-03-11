package usecase

import (
	"context"
	"database/sql"
	"strings"
	"udev21/auth/domain"

	myErrors "udev21/auth/error"
)

type userUseCase struct {
	userRepo            domain.IUserRepository
	passwordHashUseCase domain.IPasswordHashUseCase
}

func New(u domain.IUserRepository, p domain.IPasswordHashUseCase) domain.IUserUseCase {
	return &userUseCase{
		userRepo:            u,
		passwordHashUseCase: p,
	}
}

func (u *userUseCase) GetOneByID(ctx context.Context, id string) (*domain.User, error) {
	if len(id) != 36 {
		return nil, myErrors.ErrNotFound
	}

	users, err := u.userRepo.GetAllByID(ctx, []string{id})

	if len(users) == 0 {
		return nil, myErrors.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (u *userUseCase) GetOneByEmail(ctx context.Context, email string) (*domain.User, error) {
	if len(email) < 3 {
		return nil, myErrors.ErrNotFound
	}
	users, err := u.userRepo.GetAllByEmail(ctx, []string{email})
	if len(users) == 0 {
		return nil, myErrors.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &users[0], nil
}

func (u *userUseCase) GetOneByAuthJWTPayload(ctx context.Context, payload *domain.AuthJWTPayload) (*domain.User, error) {
	if payload.UserID == "" {
		return nil, myErrors.ErrNotFound
	}

	users, err := u.userRepo.GetAllByID(ctx, []string{payload.UserID})

	if len(users) == 0 {
		return nil, myErrors.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &users[0], nil
}

func (u *userUseCase) GetAllByID(ctx context.Context, userIds []string) ([]domain.User, error) {
	if len(userIds) < 1 {
		return []domain.User{}, nil
	}
	return u.userRepo.GetAllByID(ctx, userIds)
}

func (u *userUseCase) GetAllEmailStartsWith(ctx context.Context, email string, offset int64, limit uint16) ([]domain.User, error) {
	if len(email) < 3 || limit == 0 {
		return []domain.User{}, nil
	}
	// escape all % characters
	email = strings.ReplaceAll(email, "%", "\\%")

	return u.userRepo.GetAllByEmailLike(ctx, email+"%", offset, limit)
}

func (u *userUseCase) GetAllEmailLike(ctx context.Context, email string, offset int64, limit uint16) ([]domain.User, error) {
	if len(email) == 0 || limit == 0 {
		return []domain.User{}, nil
	}

	email = strings.ReplaceAll(email, "%", "\\%")

	return u.userRepo.GetAllByEmailLike(ctx, "%"+email+"%", offset, limit)
}

func (u *userUseCase) GetAllByPosition(ctx context.Context, lastPosition uint64, limit uint16) ([]domain.User, error) {
	if limit == 0 {
		return []domain.User{}, nil
	}
	return u.userRepo.GetAllByPosition(ctx, lastPosition, limit)
}

func (u *userUseCase) Create(ctx context.Context, user *domain.UserCreateInput, owner *domain.User) (*domain.User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	email := user.Email
	userExists, err := u.userRepo.GetOneByEmail(ctx, email)
	if err == nil && userExists != nil {
		return nil, myErrors.ErrEmailAlreadyExists
	} else if err != nil && err == sql.ErrNoRows {
		input := domain.User{
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Additional: user.Additional,
			Password:   u.passwordHashUseCase.Hash(user.Password),
		}
		if owner != nil {
			input.CreatedBy = &owner.ID
		}

		user, err := u.userRepo.Create(ctx, &input)

		if err != nil {
			return nil, myErrors.ErrSomethingWentWrong
		}

		return user, nil
	}
	return nil, myErrors.ErrSomethingWentWrong
}

func (u *userUseCase) Update(ctx context.Context, input *domain.UserUpdateWithoutPasswordInput) (*domain.User, error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	userExists, err := u.userRepo.GetOneByID(ctx, input.ID)

	if err != nil || userExists == nil {
		return nil, err
	}

	return u.userRepo.Update(ctx, input)
}
