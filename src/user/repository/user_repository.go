package repository

import "test-mkp/src/user/entity"

type UserRepository interface {
	FindUserByPhoneOrEmail(identity string) (*entity.User, error)
	RegisterUser(u entity.User) error
}
