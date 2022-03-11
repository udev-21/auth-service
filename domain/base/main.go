package base

import (
	"time"
	"udev21/auth/config"
	myErrors "udev21/auth/error"
)

type Model struct {
	ID        string     `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy string     `json:"created_by" db:"created_by"`
	UpdatedBy *string    `json:"updated_by" db:"updated_by"`
	DeletedBy *string    `json:"deleted_by" db:"deleted_by"`
	Position  int64      `json:"-" db:"position"`
}

func (m Model) ValidateID() error {
	if len(m.ID) < config.DBTableIDKeyLength {
		return myErrors.ErrInvalidID
	}

	return nil
}
