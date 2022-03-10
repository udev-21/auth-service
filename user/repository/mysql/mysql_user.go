package mysql

import (
	"context"
	"udev21/auth/domain"

	myErrors "udev21/auth/error"

	sqb "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
)

type mysqlUserRepository struct {
	db *sqlx.DB
}

func NewMysqlUserRepository(db *sqlx.DB) *mysqlUserRepository {
	return &mysqlUserRepository{db: db}
}

func (u *mysqlUserRepository) GetAllByByID(ctx context.Context, userIds []string) ([]domain.User, error) {
	if len(userIds) == 0 {
		return nil, myErrors.ErrNotFound
	}
	var user []domain.User = make([]domain.User, 0)
	sql, args, err := sqb.Select("*").From(domain.User{}.GetMysqlTableName()).Where(sqb.Eq{"id": userIds}).Where(sqb.Eq{"deleted_at": nil}).ToSql()

	if err != nil {
		return nil, err
	}
	err = u.db.SelectContext(ctx, &user, sql, args...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *mysqlUserRepository) GetAllByByEmail(ctx context.Context, userEmails []string) ([]domain.User, error) {
	if len(userEmails) == 0 {
		return nil, myErrors.ErrNotFound
	}
	var user []domain.User = make([]domain.User, 0)
	sql, args, err := sqb.Select("*").From(domain.User{}.GetMysqlTableName()).
		Where(sqb.Eq{"email": userEmails}).
		Where(sqb.Eq{"deleted_at": nil}).ToSql()

	if err != nil {
		return nil, err
	}
	err = u.db.SelectContext(ctx, &user, sql, args...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *mysqlUserRepository) GetAllByEmailLike(ctx context.Context, email string, offset int64, limit uint16) ([]domain.User, error) {
	if len(email) < 3 {
		return nil, myErrors.ErrNotFound
	} else if limit == 0 {
		return []domain.User{}, nil
	}

	var user []domain.User = make([]domain.User, 0)

	sql, args, err := sqb.Select("*").From(domain.User{}.GetMysqlTableName()).
		Where(sqb.Like{"email": email}).
		Where(sqb.Eq{"deleted_at": nil}).
		OrderBy("email").
		Offset(uint64(offset)).
		Limit(uint64(limit)).ToSql()

	if err != nil {
		return nil, err
	}

	err = u.db.SelectContext(ctx, &user, sql, args...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *mysqlUserRepository) GetAllByPosition(ctx context.Context, lastPosition uint64, limit uint16) ([]domain.User, error) {
	if limit == 0 {
		return []domain.User{}, nil
	}

	var user []domain.User = make([]domain.User, 0)

	sql, args, err := sqb.Select("*").From(domain.User{}.GetMysqlTableName()).
		Where(sqb.Gt{"position": lastPosition}).
		Where(sqb.Eq{"deleted_at": nil}).
		OrderBy("position").
		Limit(uint64(limit)).ToSql()

	if err != nil {
		return nil, err
	}

	err = u.db.SelectContext(ctx, &user, sql, args...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *mysqlUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	res, err := sqb.Insert(domain.User{}.GetMysqlTableName()).
		Columns("email", "first_name", "last_name", "password", "additional").
		Values(user.Email, user.FirstName, user.LastName, user.Password, user.Additional).
		RunWith(u.db).ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	if lastid, err := res.LastInsertId(); err != nil {
		return nil, myErrors.ErrSomethingWentWrong
	} else {
		return u.GetOneByPosition(ctx, uint64(lastid))
	}
}

func (u *mysqlUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	res, err := sqb.Update(domain.User{}.GetMysqlTableName()).
		Set("email", user.Email).
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("additional", user.Additional).
		Where(sqb.Eq{"id": user.ID}).
		RunWith(u.db).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	if rows, err := res.RowsAffected(); err != nil {
		return nil, err
	} else if rows == 0 {
		return nil, myErrors.ErrNotUpdated
	}
	return u.GetOneByPosition(ctx, uint64(user.Position))
}

func (u *mysqlUserRepository) UpdatePassword(ctx context.Context, user *domain.User) error {
	res, err := sqb.Update(domain.User{}.GetMysqlTableName()).
		Set("password", user.Password).
		Where(sqb.Eq{"id": user.ID}).
		RunWith(u.db).
		ExecContext(ctx)
	if err != nil {
		return myErrors.ErrSomethingWentWrong
	}
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return myErrors.ErrNotUpdated
	}
	return nil
}

func (u *mysqlUserRepository) GetOneByPosition(ctx context.Context, position uint64) (*domain.User, error) {
	var user domain.User

	sql, args, err := sqb.Select("*").From(domain.User{}.GetMysqlTableName()).
		Where(sqb.Eq{"position": position}).
		Where(sqb.Eq{"deleted_at": nil}).
		Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	err = u.db.GetContext(ctx, &user, sql, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *mysqlUserRepository) GetOneByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	sql, args, err := sqb.Select("*").From(domain.User{}.GetMysqlTableName()).
		Where(sqb.Eq{"email": email}).
		Where(sqb.Eq{"deleted_at": nil}).
		Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	err = u.db.GetContext(ctx, &user, sql, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
