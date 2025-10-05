package user

import (
	"database/sql"
	"errors"
	"fmt"
	"real-time-forum/architecture/models"
	"strings"
)

func (u *UserService) GetByNicknameOrEmail(field string) (*models.User, error) {
	// Email
	if strings.Contains(field, "@") {
		user := &models.User{Email: field}
		if err := ValidateEmail(user); err != nil {
			return nil, ErrInvalidEmail
		}

		usr, err := u.repo.GetByEmail(field)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, fmt.Errorf("u.repo.GetByEmail: %w", err)
		}
		return usr, nil
	}

	// Nickname
	user := &models.User{Nickname: field}
	if err := ValidateNickname(user); err != nil {
		return nil, ErrInvalidNickname
	}

	usr, err := u.repo.GetByNickname(field)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("u.repo.GetByNickname: %w", err)
	}
	return usr, nil
}
