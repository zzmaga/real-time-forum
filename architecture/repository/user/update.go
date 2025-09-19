package user

import (
	"database/sql"
	"fmt"
	"real-time-forum/architecture/models"
)

func (u *UserRepo) Update(user *models.User) error {
	// формируем SQL-запрос
	result, err := u.DB.Exec(`
		UPDATE users 
		SET nickname = ?, email = ?, password = ?, updated_at = ? 
		WHERE id = ?`,
		user.Nickname, user.Email, user.Password, user.UpdatedAt.Format(models.TimeFormat), user.ID,
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
