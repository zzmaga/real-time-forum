package user

import "real-time-forum/architecture/models"

type UserService struct {
	repo models.UserRepo
}

func NewUserService(repo models.UserRepo) *UserService {
	return &UserService{repo: repo}
}
