package private_message

import "real-time-forum/architecture/models"

func (s *PrivateMessageService) GetLastMessageBetweenUsers(userID1, userID2 int64) (*models.PrivateMessage, error) {
	return s.repo.GetLastMessageBetweenUsers(userID1, userID2)
}
