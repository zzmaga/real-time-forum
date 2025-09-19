package category

import (
	"errors"
	"fmt"
	rcategory "real-time-forum/architecture/repository/category"
)

func (c *CategoryService) DeleteFromPost(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid category ID: %d", id)
	}

	err := c.repo.DeleteByID(id)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, rcategory.ErrNotFound):
		return ErrNotFound
	}
	return fmt.Errorf("c.repo.DeleteByID: %w", err)
}
