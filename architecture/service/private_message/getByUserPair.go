package private_message

import "real-time-forum/architecture/models"

func (s *PrivateMessageService) GetByUserPair(userID1, userID2 int64, offset, limit int64) ([]*models.PrivateMessage, error) {
	return s.repo.GetByUserPair(userID1, userID2, offset, limit)
}
