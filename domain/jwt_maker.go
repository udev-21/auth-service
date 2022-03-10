package domain

import (
	"context"
)

type IJWTMakerUseCase interface {
	VerifyToken(ctx context.Context, token string) (*AuthJWTPayload, error)
	CreateToken(ctx context.Context, userId string) (*AuthJWT, error)
}
