package model

import "time"

type User struct {
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}

func (u *User) String() string {
	return "Created at: " + u.CreatedAt.String() +
		", Display name: " + u.DisplayName +
		", Email: " + u.Email
}

type UserList map[string]User

type UserStore struct {
	Increment int64    `json:"increment"`
	List      UserList `json:"list"`
}
