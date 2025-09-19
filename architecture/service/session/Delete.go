package session

import (
	"errors"
	"fmt"
	rsession "real-time-forum/architecture/repository/session"
)

func (s *SessionService) Delete(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid session ID: %d", id)
	}

	err := s.repo.Delete(id)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, rsession.ErrNotFound):
		return ErrNotFound
	}
	return fmt.Errorf("s.repo.Delete: %w", err)
}
