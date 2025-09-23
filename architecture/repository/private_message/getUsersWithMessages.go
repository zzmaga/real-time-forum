package private_message

import (
	"real-time-forum/architecture/models"
)

func (r *PrivateMessageRepo) GetUsersWithMessages(userID int64) ([]*models.User, error) {
	query := `
		SELECT DISTINCT u.id, u.nickname, u.email, u.first_name, u.last_name, u.age, u.gender, u.created_at, u.updated_at
		FROM users u
		WHERE u.id IN (
			SELECT DISTINCT sender_id FROM private_messages WHERE recipient_id = ?
			UNION
			SELECT DISTINCT recipient_id FROM private_messages WHERE sender_id = ?
		)
		AND u.id != ?
		ORDER BY u.nickname
	`
	rows, err := r.db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName,
			&user.Age, &user.Gender, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
