package category

import (
	"database/sql"
	"fmt"
)

func (c *CategoryRepo) DeleteByID(id int64) error {
	result, err := c.db.Exec(`DELETE FROM categories WHERE id = ?`, id)
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
