package user

import (
	"database/sql"
	"fmt"
	"real-time-forum/architecture/models"
)

func (u *UserRepo) Update(user *models.User) error {
	result, err := u.DB.Exec(`
		UPDATE users 
		SET nickname = ?, email = ?, password = ?, first_name = ?, last_name = ?, age = ?, gender = ?, updated_at = ? 
		WHERE id = ?`,
		user.Nickname, user.Email, user.Password, user.FirstName, user.LastName, user.Age, user.Gender,
		user.UpdatedAt.Format(models.TimeFormat), user.ID,
	)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}

	// проверяем, реально ли обновилась строка
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // если пользователь с таким ID не найден
	}

	return nil
}
