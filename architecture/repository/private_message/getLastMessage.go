package private_message

import (
	"database/sql"
	"real-time-forum/architecture/models"
)

func (r *PrivateMessageRepo) GetLastMessageBetweenUsers(userID1, userID2 int64) (*models.PrivateMessage, error) {
	query := `
		SELECT pm.id, pm.sender_id, pm.recipient_id, pm.content, pm.created_at,
		       u1.nickname as sender_nickname, u2.nickname as recipient_nickname
		FROM private_messages pm
		JOIN users u1 ON pm.sender_id = u1.id
		JOIN users u2 ON pm.recipient_id = u2.id
		WHERE (pm.sender_id = ? AND pm.recipient_id = ?) 
		   OR (pm.sender_id = ? AND pm.recipient_id = ?)
		ORDER BY pm.created_at DESC
		LIMIT 1
	`

	msg := &models.PrivateMessage{}
	err := r.db.QueryRow(query, userID1, userID2, userID2, userID1).Scan(
		&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.Content, &msg.CreatedAt,
		&msg.SenderNickname, &msg.RecipientNickname,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return msg, nil
}
