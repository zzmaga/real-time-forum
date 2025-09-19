package user

import (
	"errors"
	"fmt"
	"real-time-forum/architecture/models"
	"time"
)

func (u *UserService) Create(user *models.User) (int64, error) {
	if err := user.ValidateNickname(); err != nil {
		return -1, ErrInvalidNickname
	} else if err := user.ValidateEmail(); err != nil {
		return -1, ErrInvalidEmail
	}

	err := user.HashPassword()
	if err != nil {
		return -1, fmt.Errorf("user.HashPassword: %w", err)
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	userId, err := u.repo.Create(user)

	switch {
	case err == nil:
		return userId, nil
	case errors.Is(err, ErrExistEmail):
		return -1, ErrExistEmail
	case errors.Is(err, ErrExistNickname):
		return -1, ErrExistNickname
	case errors.Is(err, ErrWrongLengthEmail):
		return -1, ErrWrongLengthEmail
	case errors.Is(err, ErrWrongLengthNickname):
		return -1, ErrWrongLengthNickname
	}
	return -1, fmt.Errorf("u.repo.Create: %w", err)
}
