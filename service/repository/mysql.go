package repository

import (
	"context"
	"udev21/auth/domain"

	myErrors "udev21/auth/error"

	sqb "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type serviceRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) domain.IServiceRepository {
	return &serviceRepository{
		db: db,
	}
}

func (s *serviceRepository) GetServiceByID(ctx context.Context, id string) (*domain.Service, error) {
	var service domain.Service

	sql, args, err := sqb.Select("*").From(service.GetMysqlTableName()).Where(sqb.Eq{"id": id}).Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	err = s.db.GetContext(ctx, &service, sql, args...)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (s *serviceRepository) GetServiceByPosition(ctx context.Context, position uint64) (*domain.Service, error) {
	var service domain.Service

	sql, args, err := sqb.Select("*").From(service.GetMysqlTableName()).Where(sqb.Eq{"position": position}).Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	err = s.db.GetContext(ctx, &service, sql, args...)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (s *serviceRepository) Create(ctx context.Context, service *domain.ServiceCreateInput) (*domain.Service, error) {
	if service.Validate() != nil || service.ValidateOwnerID() != nil {
		return nil, myErrors.ErrInvalidInput
	}

	sql, args, err := sqb.Insert(service.GetMysqlTableName()).Columns("owner_id", "name", "description").Values(service.OwnerID, service.Name, service.Description).ToSql()

	if err != nil {
		return nil, err
	}

	res, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	} else if position, err := res.LastInsertId(); err != nil {
		return nil, myErrors.ErrNotCreated
	} else {
		return s.GetServiceByPosition(ctx, uint64(position))
	}
}

func (s *serviceRepository) GetAllServiceByOwnerID(ctx context.Context, ownerID string) ([]domain.Service, error) {
	var services []domain.Service

	sql, args, err := sqb.Select("*").From(domain.Service{}.GetMysqlTableName()).Where(sqb.Eq{"owner_id": ownerID}).ToSql()

	if err != nil {
		return nil, err
	}

	err = s.db.SelectContext(ctx, &services, sql, args...)

	if err != nil {
		return nil, err
	}

	return services, nil
}
