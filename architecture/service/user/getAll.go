package user

import "real-time-forum/architecture/models"

func (s *UserService) GetAll() ([]*models.User, error) {
	return s.repo.GetAll()
}
