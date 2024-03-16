package middleware

import (
	"test-mkp/config"
	"test-mkp/exception"
	"test-mkp/src/user/entity"
	"test-mkp/src/user/enum"
	"test-mkp/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

var MemberAuthMiddleware *AuthMiddleware

type AuthMiddleware struct {
	db  *sqlx.DB
	jwt config.JwtConfig
}

func NewAuthMiddleware(db *sqlx.DB, jwt config.JwtConfig) *AuthMiddleware {
	return &AuthMiddleware{
		db,
		jwt,
	}
}

func (m *AuthMiddleware) AuthMiddleware(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	if len(headers["Authorization"]) == 0 {
		exception.PanicIfNeeded(exception.UnauthenticatedError{
			Message: "unauthenticated",
		})
	}
	authorization := headers["Authorization"][0]

	if authorization == "" {
		exception.PanicIfNeeded(exception.UnauthenticatedError{
			Message: "unauthenticated",
		})
	}

	parseJwt, err := m.jwt.Parse(authorization)
	if err != nil {
		exception.PanicIfNeeded(exception.UnauthenticatedError{
			Message: "unauthenticated",
		})
	}
	xUser := parseJwt["id"].(string)

	user, err := m.getAuthObject(xUser)
	if err != nil || user == nil {
		exception.PanicIfNeeded(exception.UnauthenticatedError{
			Message: "unauthenticated",
		})
	}

	if user.Status == enum.SUSPEND {
		exception.PanicIfNeeded(exception.UnauthenticatedError{
			Message: "suspended",
		})
	}

	if user.Status == enum.INACTIVE {
		exception.PanicIfNeeded(exception.UnauthenticatedError{
			Message: "inactive",
		})
	}

	c.Locals("userId", xUser)
	return c.Next()
}

func (m *AuthMiddleware) getAuthObject(id string) (*entity.User, error) {
	ctx, cancel := config.NewPostgreSqlxContext()
	defer cancel()

	sql, args, err := sq.Select(
		"id",
		"name",
		"email",
		"phone",
		"status",
	).
		From("users").
		Where(sq.Expr("deleted_at is NULL")).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var member entity.User
	if err := m.db.GetContext(ctx, &member, sql, args...); err != nil {
		if utils.DbError(err).IsNotFound() {
			return nil, nil
		}

		return nil, err
	}

	return &member, nil
}
