package v1

import "refactoring/internal/model"

type Config interface {
	RunAddress() string
	String() string
}

type Storage interface {
	GetAllUsers() (allUsers *model.UserStore, err error)
	CreateUser(newUser *model.User) (id int64, err error)
	GetUser(id int64) (user *model.User, err error)
	UpdateUser(id int64, newDisplayName string) (err error)
	DeleteUser(id int64) (err error)
}
