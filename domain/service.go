package domain

import (
	"strings"
	"udev21/auth/domain/base"
	myErrors "udev21/auth/error"
)

type Service struct {
	base.Model
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	UserId      string  `json:"user_id" db:"user_id"`
}

func (s Service) GetMysqlTableName() string {
	return "services"
}

func (s *Service) ValidateName() error {
	s.Name = strings.TrimSpace(s.Name)
	if len(s.Name) < 1 {
		return myErrors.ErrInvalidValue
	}
	return nil
}

func (s *Service) ValidateDescription() error {
	if s.Description != nil {
		*s.Description = strings.TrimSpace(*s.Description)
		if len(*s.Description) < 1 {
			return myErrors.ErrInvalidValue
		}
	}
	return nil
}

func (s *Service) ValidateUserID() error {
	if s.Description != nil {
		*s.Description = strings.TrimSpace(*s.Description)
		if len(*s.Description) < 1 {
			return myErrors.ErrInvalidValue
		}
	}
	return nil
}
