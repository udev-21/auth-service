package domain

import (
	"context"
)

type IAuthUseCase interface {
	Login(ctx context.Context, input UserLoginInput) (*AuthJWT, error)
	Register(ctx context.Context, input UserCreateInput) (*AuthJWT, error)
	RefreshToken(ctx context.Context, input AuthJWT) (*AuthJWT, error)
}
