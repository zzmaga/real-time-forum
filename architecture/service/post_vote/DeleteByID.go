package post_vote

import (
	"database/sql"
	"errors"
	"fmt"
)

func (p *PostVoteService) DeleteByID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid postvote ID: %d", id)
	}

	err := p.repo.DeleteByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("p.repo.DeleteByID: %w", err)
	}

	return nil
}
