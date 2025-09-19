package user

import (
	"errors"
	"fmt"
	ruser "real-time-forum/architecture/repository/user"
)

func (u *UserService) DeleteByID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID: %d", id)
	}

	err := u.repo.DeleteByID(id)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ruser.ErrNotFound):
		return ErrNotFound
	}
	return fmt.Errorf("u.repo.DeleteByID: %w", err)
}
