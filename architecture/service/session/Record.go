package session

import (
	"errors"
	"fmt"
	"time"

	"real-time-forum/architecture/models"

	rsession "real-time-forum/architecture/repository/session"

	uuid "github.com/satori/go.uuid"
)

func (s *SessionService) Record(userId int64) (*models.Session, error) {
	uid := uuid.NewV4()
	session := &models.Session{
		Uuid:      uid.String(),
		UserID:    userId,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(models.SessionExpiredAfter),
	}

	_, err := s.repo.Create(session)
	switch {
	case err == nil:
		return session, nil
	case errors.Is(err, rsession.ErrSessionExists):
		err := s.repo.UpdateByUserId(session.UserID, session)
		if err != nil {
			return nil, fmt.Errorf("s.repo.UpdateByUserId: %w", err)
		}
		return session, nil
	default:
		return nil, fmt.Errorf("s.repo.Create: %w", err)
	}
}
