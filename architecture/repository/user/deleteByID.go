package user

import (
	"database/sql"
	"fmt"
)

func (u *UserRepo) DeleteByID(id int64) error {
	result, err := u.DB.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("DeleteByID: %w", err)
	}

	// проверим, был ли реально удалён пользователь
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteByID: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // никого не удалили → такого id нет
	}

	return nil
}
