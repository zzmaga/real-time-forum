package user

import (
	"fmt"
	"strings"
	"time"

	"real-time-forum/architecture/models"
)

func (u *UserRepo) GetByEmail(email string) (*models.User, error) {
	row := u.DB.QueryRow(`
SELECT id, nickname, email, password, first_name, last_name, age, gender, created_at, updated_at FROM users
WHERE email = ?`, email)
	user := &models.User{}
	strCreatedAt, strUpdatedAt := "", ""

	err := row.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &user.Age, &user.Gender, &strCreatedAt, &strUpdatedAt)

	switch {
	case err == nil:
		timeCreatedAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse created_at: %w", err)
		}
		timeUpdatedAt, err := time.ParseInLocation(models.TimeFormat, strUpdatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse updated_at: %w", err)
		}
		user.CreatedAt = timeCreatedAt
		user.UpdatedAt = timeUpdatedAt
		return user, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
}
