package post_vote

import (
	"real-time-forum/architecture/models"
)

type PostVoteService struct {
	repo models.PostVoteRepo
}

func NewPostVoteService(postVote models.PostVoteRepo) *PostVoteService {
	return &PostVoteService{postVote}
}
