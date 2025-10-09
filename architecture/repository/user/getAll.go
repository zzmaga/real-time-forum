package user

import (
	"fmt"
	"time"

	"real-time-forum/architecture/models"
)

func (u *UserRepo) GetAll() ([]*models.User, error) {
	rows, err := u.DB.Query(`
		SELECT id, nickname, email, first_name, last_name, age, gender, created_at, updated_at 
		FROM users
		WHERE nickname IS NOT NULL AND nickname != ''
		ORDER BY nickname
	`)

	if err != nil {
		return nil, fmt.Errorf("u.db.Query: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var strCreatedAt, strUpdatedAt string

		err := rows.Scan(&user.ID, &user.Nickname, &user.Email,
			&user.FirstName, &user.LastName, &user.Age, &user.Gender, &strCreatedAt, &strUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
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
		users = append(users, user)
	}

	return users, nil
}
