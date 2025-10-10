package private_message

import (
	"real-time-forum/architecture/models"
)

func (r *PrivateMessageRepo) GetByUserPair(userID1, userID2 int64, offset, limit int64) ([]*models.PrivateMessage, error) {
	query := `
		SELECT pm.id, pm.sender_id, pm.recipient_id, pm.content, pm.created_at,
		       u1.nickname as sender_nickname, u2.nickname as recipient_nickname
		FROM private_messages pm
		JOIN users u1 ON pm.sender_id = u1.id
		JOIN users u2 ON pm.recipient_id = u2.id
		WHERE (pm.sender_id = ? AND pm.recipient_id = ?) 
		   OR (pm.sender_id = ? AND pm.recipient_id = ?)
		ORDER BY pm.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, userID1, userID2, userID2, userID1, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []*models.PrivateMessage
	for rows.Next() {
		msg := &models.PrivateMessage{}
		err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.Content, &msg.CreatedAt,
			&msg.SenderNickname, &msg.RecipientNickname,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if messages == nil {
		// Js воспринимает нил как null из за чего возникает ошибка на фронте
		return []*models.PrivateMessage{}, nil
	}
	return messages, nil
}
