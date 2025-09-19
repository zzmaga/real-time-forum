package category

import (
	"fmt"
	"strings"
	"time"

	"real-time-forum/architecture/models"
)

func (c *CategoryRepo) GetByID(id int64) (*models.Category, error) {
	row := c.db.QueryRow(`
SELECT id, name, created_at FROM categories
WHERE id = ?`, id)
	category := &models.Category{}
	strCreatedAt := ""
	err := row.Scan(&category.Id, &category.Name, &strCreatedAt)
	switch {
	case err == nil:
		timeCreatedAt, err := time.ParseInLocation(models.TimeFormat, strCreatedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("time.Parse: %w", err)
		}
		category.CreatedAt = timeCreatedAt
		return category, nil
	case strings.HasPrefix(err.Error(), "sql: no rows in result set"):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
}
