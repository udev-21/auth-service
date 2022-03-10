package domain

import (
	"time"

	myErrors "udev21/auth/error"
)

type AuthJWT struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthJWTPayload struct {
	UserID    string    `json:"uid"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (payload *AuthJWTPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return myErrors.ErrAuthTokenExpired
	}
	return nil
}

func NewAuthJWTPayload(userId string, duration time.Duration) *AuthJWTPayload {
	return &AuthJWTPayload{
		UserID:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}
