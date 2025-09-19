package user

import (
	"fmt"
	"strings"

	"real-time-forum/architecture/models"
)

func (u *UserRepo) Create(user *models.User) (int64, error) {
	strCreatedAt := user.CreatedAt.Format(models.TimeFormat)
	strUpdatedAt := user.UpdatedAt.Format(models.TimeFormat)
	row := u.DB.QueryRow(`
INSERT INTO users (nickname, email, password, first_name, last_name, age, gender, created_at, updated_at) VALUES
(?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		user.Nickname, user.Email, user.Password, user.FirstName, user.LastName, user.Age, user.Gender,
		strCreatedAt, strUpdatedAt)

	err := row.Scan(&user.ID)
	switch {
	case err == nil:
		return user.ID, nil
	case strings.HasPrefix(err.Error(), "UNIQUE constraint failed"):
		switch {
		case strings.Contains(err.Error(), "nickname"):
			return -1, ErrExistNickname
		case strings.Contains(err.Error(), "email"):
			return -1, ErrExistEmail
		}
	case strings.HasPrefix(err.Error(), "CHECK constraint failed"):
		switch {
		case strings.Contains(err.Error(), "LENGTH(nickname)"):
			return -1, ErrWrongLengthNickname
		case strings.Contains(err.Error(), "LENGTH(email)"):
			return -1, ErrWrongLengthEmail
		}
	}
	return -1, fmt.Errorf("row.Scan: %w", err)
}
