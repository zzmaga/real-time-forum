package post_vote

import (
	"database/sql"
	"fmt"
)

func (p *PostVoteRepo) DeleteByID(id int64) error {
	result, err := p.db.Exec(`DELETE FROM post_votes WHERE id = ?`, id)
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
