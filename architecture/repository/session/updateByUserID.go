package session

import (
	"fmt"

	"real-time-forum/architecture/models"
)

func (s *SessionRepo) UpdateByUserId(userId int64, session *models.Session) error {
	strExpiredAt := session.ExpiredAt.Format(models.TimeFormat)
	row := s.db.QueryRow(`
UPDATE sessions 
SET uuid = ?, expired_at = ?
WHERE user_id = ?
RETURNING id`, session.Uuid, strExpiredAt, session.UserID)

	err := row.Scan(&session.ID)
	switch {
	case err == nil:
		return nil
	}
	return fmt.Errorf("row.Scan: %w", err)
}
