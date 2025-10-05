package user

import (
	"database/sql"
	"errors"
	"fmt"
	"real-time-forum/architecture/models"
)

func (u *UserService) GetByID(id int64) (*models.User, error) {
	usr, err := u.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("u.repo.GetByID: %w", err)
	}
	return usr, nil
}
