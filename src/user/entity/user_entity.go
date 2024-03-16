package entity

import (
	"test-mkp/src/user/enum"
	"time"
)

type User struct {
	ID       string          `json:"id" db:"id"`
	Name     *string         `json:"name" db:"name"`
	Email    *string         `json:"email" db:"email"`
	Phone    *string         `json:"phone" db:"phone"`
	Photo    *string         `json:"photo" db:"photo"`
	Birthday *time.Time      `json:"birthday" db:"birthday"`
	Status   enum.UserStatus `json:"status" db:"status"`
	Password string          `db:"password"`
	Default
}
