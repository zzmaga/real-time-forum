package models

import "time"

type PrivateMessage struct {
	ID          int64
	SenderID    int64
	RecipientID int64
	Content     string
	CreatedAt   time.Time

	// Additional fields for display
	SenderNickname    string
	RecipientNickname string
}

type PrivateMessageRepo interface {
	Create(message *PrivateMessage) (int64, error)
	GetByUserPair(userID1, userID2 int64, offset, limit int64) ([]*PrivateMessage, error)
	GetLastMessageBetweenUsers(userID1, userID2 int64) (*PrivateMessage, error)
	GetUsersWithMessages(userID int64) ([]*User, error)
}

type PrivateMessageService interface {
	Create(message *PrivateMessage) (int64, error)
	GetByUserPair(userID1, userID2 int64, offset, limit int64) ([]*PrivateMessage, error)
	GetLastMessageBetweenUsers(userID1, userID2 int64) (*PrivateMessage, error)
	GetUsersWithMessages(userID int64) ([]*User, error)
}
