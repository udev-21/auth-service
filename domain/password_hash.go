package domain

type IPasswordHashUseCase interface {
	Hash(password string) string
	Compare(hashedPassword, password string) bool
}
