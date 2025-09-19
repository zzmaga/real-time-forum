package category

import (
	"database/sql"
	"fmt"
	"real-time-forum/architecture/models"
)

func (c *CategoryRepo) Update(category *models.Category) error {
	result, err := c.db.Exec(`
		UPDATE categories 
		SET name = ? 
		WHERE id = ?`,
		category.Name, category.Id,
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
		return sql.ErrNoRows // если категория с таким ID не найдена
	}

	return nil
}
