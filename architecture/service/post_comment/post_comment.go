package post_comment

import (
	"real-time-forum/architecture/models"
)

type PostCommentService struct {
	repo models.PostCommentRepo
}

func NewPostCommentService(repo models.PostCommentRepo) *PostCommentService {
	return &PostCommentService{repo}
}
