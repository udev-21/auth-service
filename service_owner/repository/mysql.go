package repository

import (
	"context"
	"udev21/auth/domain"

	sqb "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type serviceOwnerRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) domain.IServiceOwnerRepository {
	return &serviceOwnerRepository{
		db: db,
	}
}

func (r *serviceOwnerRepository) GetServiceOwnerByID(ctx context.Context, id string) (*domain.ServiceOwner, error) {
	var serviceOwner domain.ServiceOwner

	sql, args, err := sqb.Select("id", "user_id").From(serviceOwner.GetMysqlTableName()).Where(sqb.Eq{"user_id": id}).Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, &serviceOwner, sql, args...)
	if err != nil {
		return nil, err
	}

	return &serviceOwner, nil
}