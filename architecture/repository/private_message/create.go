package private_message

import (
	"real-time-forum/architecture/models"
	"time"
)

func (r *PrivateMessageRepo) Create(message *models.PrivateMessage) (int64, error) {
	query := `
		INSERT INTO private_messages (sender_id, recipient_id, content, created_at)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, message.SenderID, message.RecipientID, message.Content, time.Now())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
