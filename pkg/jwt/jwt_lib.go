package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtProvider interface {
	Generate(claims JwtClaims) (string, error)
	ParseJwt(jwtToken string) (jwt.MapClaims, error)
}

type JwtProviderConfig struct {
	SigningMethod    *jwt.SigningMethodHMAC
	SigningSiganture []byte
}

type JwtClaims struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
	Id        string
}

type jwtProviderImpl struct {
	SigningMethod    *jwt.SigningMethodHMAC
	SigningSiganture []byte
}

type claimRegistration struct {
	jwt.RegisteredClaims
	Id string `json:"id"`
}

func New(config JwtProviderConfig) JwtProvider {
	return &jwtProviderImpl{
		SigningMethod:    config.SigningMethod,
		SigningSiganture: config.SigningSiganture,
	}
}

func (j *jwtProviderImpl) Generate(claims JwtClaims) (string, error) {
	c := claimRegistration{
		Id: claims.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(claims.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
		},
	}

	token := jwt.NewWithClaims(j.SigningMethod, c)
	signed, err := token.SignedString(j.SigningSiganture)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (j *jwtProviderImpl) ParseJwt(jwtToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unauthenticated")
		} else if method != j.SigningMethod {
			return nil, fmt.Errorf("unauthenticate")
		}

		return j.SigningSiganture, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("claim error %v - %v", ok, token.Valid)
	}

	return claims, nil
}
