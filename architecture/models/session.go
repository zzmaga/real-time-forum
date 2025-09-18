package models

import (
	"time"
)

type Session struct {
	ID        int64
	Uuid      string
	UserID    int64
	CreatedAt time.Time
	ExpiredAt time.Time
}

type SessionRepo interface {
	Create(session *Session) (int64, error)
	Delete(id int64) error
	GetByUuid(uuid string) (*Session, error)
	UpdateByUserId(userId int64, session *Session) error
}

type SessionService interface {
	Record(userId int64) (*Session, error)
	Delete(id int64) error
	GetByUuid(uuid string) (*Session, error)
}
