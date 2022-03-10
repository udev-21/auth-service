package usecase

import (
	"context"
	"database/sql"
	"time"
	"udev21/auth/domain"
	myErrors "udev21/auth/error"
)

var passwordSalt = []byte("asdfasdf")
var tokenTimeout = time.Minute * 15
var refreshTokenTimeout = time.Hour * 24 * 7

type authUseCase struct {
	userRepo            domain.IUserRepository
	jwtMakerUseCase     domain.IJWTMakerUseCase
	passwordHashUseCase domain.IPasswordHashUseCase
}

func NewAuthUseCase(userRepo domain.IUserRepository, jwtMakerUseCase domain.IJWTMakerUseCase, passwordHashUseCase domain.IPasswordHashUseCase) domain.IAuthUseCase {
	return &authUseCase{
		userRepo:            userRepo,
		jwtMakerUseCase:     jwtMakerUseCase,
		passwordHashUseCase: passwordHashUseCase,
	}
}

func (u *authUseCase) Login(ctx context.Context, input domain.UserInput) (*domain.AuthJWT, error) {
	if err := input.ValidateEmail(); err != nil {
		return nil, myErrors.ErrNotFound
	}

	user, err := u.userRepo.GetOneByEmail(ctx, input.Email)
	if err != nil {
		return nil, myErrors.ErrPasswordOrEmailIncorrect
	}

	if !u.passwordHashUseCase.Compare(user.Password, input.Password) {
		return nil, myErrors.ErrPasswordOrEmailIncorrect
	}

	token, err := u.jwtMakerUseCase.CreateToken(ctx, user.ID)

	if err != nil {
		return nil, myErrors.ErrSomethingWentWrong
	}

	return token, nil
}

func (u *authUseCase) Register(ctx context.Context, input domain.UserInput) (*domain.AuthJWT, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	email := input.Email
	userExists, err := u.userRepo.GetOneByEmail(ctx, email)
	if err == nil && userExists != nil {
		return nil, myErrors.ErrEmailAlreadyExists
	} else if err != nil && err == sql.ErrNoRows {

		user, err := u.userRepo.Create(ctx, &domain.User{
			Email:      input.Email,
			FirstName:  input.FirstName,
			LastName:   input.LastName,
			Additional: input.Additional,
			Password:   u.passwordHashUseCase.Hash(input.Password),
		})
		if err != nil {
			return nil, myErrors.ErrSomethingWentWrong
		}
		token, err := u.jwtMakerUseCase.CreateToken(ctx, user.ID)
		if err != nil {
			panic(err)
		}

		return token, nil
	}
	return nil, myErrors.ErrSomethingWentWrong

}

func (u *authUseCase) RefreshToken(ctx context.Context, input domain.AuthJWT) (*domain.AuthJWT, error) {
	token, err := u.jwtMakerUseCase.VerifyToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, myErrors.ErrNotFound
	} else if len(token.UserID) != 36 {
		return nil, myErrors.ErrNotFound
	}

	return u.jwtMakerUseCase.CreateToken(ctx, token.UserID)
}
