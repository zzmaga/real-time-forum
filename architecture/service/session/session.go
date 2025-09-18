package session

import "real-time-forum/architecture/models"

type SessionService struct {
	repo models.SessionRepo
}

func NewSessionService(repo models.SessionRepo) *SessionService {
	return &SessionService{repo}
}
