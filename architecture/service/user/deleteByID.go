package user

import (
	"database/sql"
	"errors"
	"fmt"
)

func (u *UserService) DeleteByID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID: %d", id)
	}

	err := u.repo.DeleteByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("u.repo.DeleteByID: %w", err)
	}

	return nil
}
