package user

import (
	"fmt"
	"time"

	"real-time-forum/architecture/models"
)

func (u *UserRepo) GetByNickname(nickname string) (*models.User, error) {
	row := u.DB.QueryRow(`
SELECT id, nickname, email, password, first_name, last_name, age, gender, created_at, updated_at 
FROM users
WHERE nickname = ?`, nickname)

	user := &models.User{}
	var strCreatedAt, strUpdatedAt string

	err := row.Scan(&user.ID, &user.Nickname, &user.Email, &user.Password,
		&user.FirstName, &user.LastName, &user.Age, &user.Gender, &strCreatedAt, &strUpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", err)
	}

	createdAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse created_at: %w", err)
	}
	updatedAt, err := time.ParseInLocation(models.TimeFormat, strUpdatedAt, time.Local)
	if err != nil {
		return nil, fmt.Errorf("time.Parse updated_at: %w", err)
	}

	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	return user, nil
}
