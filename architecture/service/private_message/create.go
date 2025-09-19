package private_message

import "real-time-forum/architecture/models"

func (s *PrivateMessageService) Create(message *models.PrivateMessage) (int64, error) {
	return s.repo.Create(message)
}
