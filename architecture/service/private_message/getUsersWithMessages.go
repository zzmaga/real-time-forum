package private_message

import "real-time-forum/architecture/models"

func (s *PrivateMessageService) GetUsersWithMessages(userID int64) ([]*models.User, error) {
	return s.repo.GetUsersWithMessages(userID)
}
