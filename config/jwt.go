package config

import (
	jwt_lib "test-mkp/pkg/jwt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig interface {
	Generate(id string) (string, error)
	Parse(token string) (jwt.MapClaims, error)
}

type jwtConfigImpl struct {
	provider jwt_lib.JwtProvider
}

func NewJwtConfig(config Config) JwtConfig {
	provider := jwt_lib.New(jwt_lib.JwtProviderConfig{
		SigningMethod:    jwt.SigningMethodHS256,
		SigningSiganture: []byte(config.Get("APP_SECRET")),
	})

	return &jwtConfigImpl{
		provider: provider,
	}
}

func (j *jwtConfigImpl) Generate(id string) (string, error) {
	claims := jwt_lib.JwtClaims{
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().AddDate(1, 0, 0),
		Id:        id,
	}

	return j.provider.Generate(claims)
}

func (j *jwtConfigImpl) Parse(token string) (jwt.MapClaims, error) {
	return j.provider.ParseJwt(token)
}
