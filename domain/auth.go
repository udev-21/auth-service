package domain

import (
	"context"
)

type IAuthUseCase interface {
	Login(ctx context.Context, input UserInput) (*AuthJWT, error)
	Register(ctx context.Context, input UserInput) (*AuthJWT, error)
	RefreshToken(ctx context.Context, input AuthJWT) (*AuthJWT, error)
}
