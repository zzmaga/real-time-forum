package private_message

import (
	"real-time-forum/architecture/models"
)

type PrivateMessageService struct {
	repo models.PrivateMessageRepo
}

func NewPrivateMessageService(repo models.PrivateMessageRepo) models.PrivateMessageService {
	return &PrivateMessageService{repo: repo}
}
