package user

import "errors"

type UserStore interface {
	GetUser(userId int64) (*User, error)
	CreateUser(*User) error
	UpdateUser(*User) error
	DeleteUser(*User) error
}

var (
	ErrNoUserFound = errors.New("no user found")
)
