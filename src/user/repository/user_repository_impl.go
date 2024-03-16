package repository

import (
	"strings"
	"test-mkp/config"
	"test-mkp/src/user/entity"
	"test-mkp/src/user/enum"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

func (r *userRepositoryImpl) FindUserByPhoneOrEmail(identity string) (*entity.User, error) {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()
	response := entity.User{}

	q := sq.Select().
		Columns("a.id", "a.name", "a.email", "a.phone", "a.photo", "a.birthday", "status", "password", "created_at", "updated_at", "deleted_at").
		From("users a").
		Where(sq.And{
			sq.Or{
				sq.Eq{
					"a.email": identity,
				},
				sq.Eq{
					"a.phone": identity,
				},
			},
			sq.Eq{
				"a.status": int(enum.ACTIVE),
			},
		}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.GetContext(ctx, &response, sql, args...); err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return nil, nil
		}

		return &response, err
	}

	return &response, nil
}

func (r *userRepositoryImpl) RegisterUser(u entity.User) error {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	statement := sq.Insert("users").
		Columns("id", "name", "email", "phone", "password", "status", "created_at", "updated_at").
		Values(u.ID, u.Name, u.Email, u.Phone, u.Password, enum.ACTIVE, u.CreatedAt, u.UpdatedAt).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := statement.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
