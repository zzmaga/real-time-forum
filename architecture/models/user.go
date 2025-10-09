package models

import "time"

type User struct {
	ID        int64
	Nickname  string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Age       int
	Gender    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserService interface {
	Create(user *User) (int64, error)
	Update(user *User) error
	DeleteByID(id int64) error

	GetByID(id int64) (*User, error)
	GetByNicknameOrEmail(field string) (*User, error)
	GetAll() ([]*User, error)
}

type UserRepo interface {
	Create(user *User) (int64, error)
	Update(user *User) error
	DeleteByID(id int64) error

	GetByID(id int64) (*User, error)
	GetByNickname(nickname string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll() ([]*User, error)
}
