package domain

import (
	"udev21/auth/config"
	"udev21/auth/domain/base"
	myErrors "udev21/auth/error"
)

type ServiceUser struct {
	base.Model
	UserID    string `json:"user_id" db:"user_id"`
	ServiceID string `json:"service_id" db:"service_id"`
}

func (s ServiceUser) GetMysqlTableName() string {
	return "service_users"
}

func (s *ServiceUser) ValidateUserID() error {
	if len(s.UserID) < config.DBTableIDKeyLength {
		return myErrors.ErrInvalidValue
	}
	return nil
}
func (s *ServiceUser) ValidateServiceID() error {
	if len(s.ServiceID) < config.DBTableIDKeyLength {
		return myErrors.ErrInvalidValue
	}
	return nil
}
