package domain

import (
	"context"
	"strings"
	"udev21/auth/domain/base"
	myErrors "udev21/auth/error"
	"udev21/auth/util"
)

type User struct {
	base.Model
	FirstName  *string `json:"first_name" db:"first_name"`
	LastName   *string `json:"last_name" db:"last_name"`
	Email      string  `json:"email" db:"email"`
	Password   string  `json:"-" db:"password"`
	Additional *string `json:"additional" db:"additional"`
}

func (u User) GetMysqlTableName() string {
	return "users"
}

type IUserUseCase interface {
	GetOneByID(ctx context.Context, id string) (*User, error)
	GetOneByEmail(ctx context.Context, email string) (*User, error)
	GetOneByAuthJWTPayload(ctx context.Context, payload *AuthJWTPayload) (*User, error)

	GetAllByID(ctx context.Context, userIds []string) ([]User, error)
	GetAllByPosition(ctx context.Context, lastPosition uint64, limit uint16) ([]User, error)

	GetAllEmailStartsWith(ctx context.Context, email string, offset int64, limit uint16) ([]User, error)
	GetAllEmailLike(ctx context.Context, email string, offset int64, limit uint16) ([]User, error)

	Create(ctx context.Context, input *UserInput) (*User, error)
	Update(ctx context.Context, input *User) (*User, error)
	// Delete(ctx context.Context, userIds []string) error
	// DeleteAll(ctx context.Context, users []User) error
}

type IUserRepository interface {
	GetAllByPosition(ctx context.Context, lastPosition uint64, limit uint16) ([]User, error)
	GetAllByByID(ctx context.Context, userIds []string) ([]User, error)
	GetAllByByEmail(ctx context.Context, userEmails []string) ([]User, error)
	GetAllByEmailLike(ctx context.Context, email string, offset int64, limit uint16) ([]User, error)
	GetOneByPosition(ctx context.Context, position uint64) (*User, error)
	GetOneByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	UpdatePassword(ctx context.Context, user *User) error
	// Delete(ctx context.Context, userIds []string) error
}

type UserInput struct {
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	PasswordConfirm *string `json:"password_confirm"`
	Additional      *string `json:"additional"`
}

func (a UserInput) Validate() error {
	return util.GetErrorIfExist(a.ValidateEmail, a.ValidatePassword, a.ValidateConfirmPassword, a.ValidateFirstName, a.ValidateLastName, a.ValidateAdditional)
}
func (a UserInput) ValidateWithoutPassword() error {
	return util.GetErrorIfExist(a.ValidateEmail, a.ValidateFirstName, a.ValidateLastName, a.ValidateAdditional)
}

func (a UserInput) ValidateEmail() error {
	return util.ValidateEmail(a.Email)
}

func (a UserInput) ValidateFirstName() error {
	if a.FirstName != nil {
		*a.FirstName = strings.TrimSpace(*a.FirstName)
		if *a.FirstName == "" {
			return myErrors.ErrInvalidInput
		}
	}
	return nil
}
func (a UserInput) ValidateAdditional() error {
	if a.Additional != nil {
		*a.Additional = strings.TrimSpace(*a.Additional)
		if *a.Additional == "" {
			return myErrors.ErrInvalidInput
		}
	}
	return nil
}

func (a UserInput) ValidateLastName() error {
	if a.LastName != nil {
		*a.LastName = strings.TrimSpace(*a.LastName)
		if *a.LastName == "" {
			return myErrors.ErrInvalidInput
		}
	}
	return nil
}

func (a UserInput) ValidatePassword() error {
	if util.ValidatePasswordStrength(a.Password) == false {
		return myErrors.ErrInvalidPassword
	}
	return nil
}

func (a UserInput) ValidateConfirmPassword() error {
	if a.Password != *a.PasswordConfirm {
		return myErrors.ErrPasswordConfirmNotMatch
	}
	return nil
}
