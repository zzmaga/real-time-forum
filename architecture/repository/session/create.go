package session

import (
	"fmt"
	"strings"

	"real-time-forum/architecture/models"
)

func (s *SessionRepo) Create(session *models.Session) (int64, error) {
	strExpiredAt := session.ExpiredAt.Format(models.TimeFormat)
	res, err := s.db.Exec(`
INSERT INTO sessions (uuid, created_at, expired_at, user_id) VALUES
(?, ?, ?, ?)`, session.Uuid, session.CreatedAt, strExpiredAt, session.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return -1, ErrSessionExists
		}
		return -1, fmt.Errorf("s.db.Exec: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("res.LastInsertId: %w", err)
	}
	session.ID = id
	return session.ID, nil
	/*
		RETURNIN id работает в постгре, в sqlite нет
		Твоя структура не сработает с sqlite, поэтому я сверху так написал
	*/
	/*row := s.db.QueryRow(`
	INSERT INTO sessions (uuid, created_at, expired_at, user_id) VALUES
	(?, ?, ?, ?) RETURNING id`, session.Uuid, session.CreatedAt, strExpiredAt, session.UserID)
		log.Println("passed")
		err := row.Scan(&session.ID)
		switch {
		case err == nil:
			log.Println("here2")
			return session.ID, nil
		case strings.HasPrefix(err.Error(), "UNIQUE constraint failed"):
			log.Println("here")
			return -1, ErrSessionExists
		}
		log.Println(err.Error())
		return -1, fmt.Errorf("row.Scan: %w", err)*/
}
