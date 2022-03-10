package usecase

import (
	"encoding/hex"
	"udev21/auth/domain"

	"golang.org/x/crypto/argon2"
)

type passwordHashUseCase struct {
	config domain.PasswordConfig
}

func New(config domain.PasswordConfig) domain.IPasswordHashUseCase {
	return &passwordHashUseCase{
		config: config,
	}
}

func (u *passwordHashUseCase) Hash(password string) string {
	return hex.EncodeToString(argon2.IDKey([]byte(password), u.config.Salt, u.config.Argon.Time, u.config.Argon.Memory, u.config.Argon.Threads, u.config.Argon.KeyLength))
}

func (u *passwordHashUseCase) Compare(hashedPassword, password string) bool {
	return u.Hash(password) == hashedPassword
}
