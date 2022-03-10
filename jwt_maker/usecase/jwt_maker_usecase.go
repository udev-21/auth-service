package usecase

import (
	"context"
	"errors"
	"fmt"
	"udev21/auth/domain"
	myErrors "udev21/auth/error"

	"github.com/golang-jwt/jwt"
)

type jwtMakerUseCase struct {
	config domain.JWTConfig
}

const minSecretKeySize = 32

func New(config domain.JWTConfig) (domain.IJWTMakerUseCase, error) {
	if len(config.SecretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &jwtMakerUseCase{
		config: config,
	}, nil
}

func (m *jwtMakerUseCase) VerifyToken(ctx context.Context, token string) (*domain.AuthJWTPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, myErrors.ErrInvalidToken
		}
		return m.config.SecretKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &domain.AuthJWTPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, myErrors.ErrExpiredToken) {
			return nil, myErrors.ErrExpiredToken
		}
		return nil, myErrors.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*domain.AuthJWTPayload)

	if !ok {
		return nil, myErrors.ErrInvalidToken
	}

	return payload, nil
}

func (m *jwtMakerUseCase) CreateToken(ctx context.Context, userId string) (*domain.AuthJWT, error) {

	accessToken, err := m.createAccesToken(userId)
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.createRefreshToken(userId)
	if err != nil {
		return nil, err
	}

	return &domain.AuthJWT{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *jwtMakerUseCase) createAccesToken(userId string) (string, error) {
	accessPayload := domain.NewAuthJWTPayload(userId, m.config.AccessTokenExpireDuration)
	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessPayload)
	return accessJwt.SignedString(m.config.SecretKey)
}

func (m *jwtMakerUseCase) createRefreshToken(userId string) (string, error) {
	refreshPayload := domain.NewAuthJWTPayload(userId, m.config.RefreshTokenExpireDuration)
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshPayload)

	return refreshJwt.SignedString(m.config.SecretKey)
}
