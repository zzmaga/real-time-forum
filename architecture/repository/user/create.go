package user

import (
	"fmt"

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
	if err != nil {
		return 0, fmt.Errorf("repo.Create: %w", err)
	}
	return user.ID, nil
}
