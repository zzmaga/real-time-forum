package post_comment_vote

import "real-time-forum/architecture/models"

type PostCommentVoteService struct {
	repo models.PostCommentVoteRepo
}

func NewPostCommentVoteService(postCommentVote models.PostCommentVoteRepo) *PostCommentVoteService {
	return &PostCommentVoteService{postCommentVote}
}
