package config

import (
	"context"
	"fmt"
	"strconv"
	"test-mkp/exception"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgreSqlxConfig struct {
	DB     *sqlx.DB
	Source string
}

func NewPostgreSqlxDatabase(config Config) PostgreSqlxConfig {
	ctx, cancel := NewPostgreSqlxContext()
	defer cancel()

	var dbConn *sqlx.DB
	host := config.Get("POSTGRE_HOST")
	user := config.Get("POSTGRE_USER")
	pass := config.Get("POSTGRE_PASS")
	db := config.Get("POSTGRE_DB")
	ssl := config.Get("POSTGRE_SSL")
	port := config.Get("POSTGRE_PORT")

	idleTime, err := strconv.Atoi(config.Get("POSTGRE_MAX_IDLE"))
	exception.PanicIfNeeded(err)

	maxLifetime, err := strconv.Atoi(config.Get("POSTGRE_MAX_LIFE_TIME"))
	exception.PanicIfNeeded(err)

	maxConn, err := strconv.Atoi(config.Get("POSTGRE_POOL_MAX"))
	exception.PanicIfNeeded(err)

	dbSource := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s dbname=%s", user, pass, host, port, ssl, db)
	dbConn, err = sqlx.ConnectContext(ctx, "pgx", dbSource)
	exception.PanicIfNeeded(err)

	dbConn.SetMaxOpenConns(maxConn)
	dbConn.SetConnMaxIdleTime(time.Duration(idleTime))
	dbConn.SetConnMaxLifetime(time.Duration(maxLifetime))

	return PostgreSqlxConfig{
		DB:     dbConn,
		Source: dbSource,
	}
}

func NewPostgreSqlxContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
