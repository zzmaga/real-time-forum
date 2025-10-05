package user

import (
	"fmt"
	"real-time-forum/architecture/models"
	"strings"
	"time"
)

func (u *UserService) Create(user *models.User) (int64, error) {
	ValidateNickname(user)
	ValidateEmail(user)
	HashPassword(user)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	userId, err := u.repo.Create(user)
	if err == nil {
		return userId, nil
	}

	errMsg := err.Error()

	switch {
	case strings.HasPrefix(errMsg, "UNIQUE constraint failed"):
		if strings.Contains(errMsg, "nickname") {
			return 0, ErrExistNickname
		}
		if strings.Contains(errMsg, "email") {
			return 0, ErrExistEmail
		}
	case strings.HasPrefix(errMsg, "CHECK constraint failed"):
		if strings.Contains(errMsg, "LENGTH(nickname)") {
			return 0, ErrWrongLengthNickname
		}
		if strings.Contains(errMsg, "LENGTH(email)") {
			return 0, ErrWrongLengthEmail
		}
	}

	return 0, fmt.Errorf("service.Create: %w", err) // Undefined err
}
