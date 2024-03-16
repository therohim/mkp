package service

import "test-mkp/src/user/model"

type UserService interface {
	Register(request model.RegisterRequest) error
	Login(request model.LoginRequest) (model.LoginResponse, error)
}
