package service

import (
	"real-time-forum/architecture/models"
	"real-time-forum/architecture/repository"
	"real-time-forum/architecture/service/category"
	"real-time-forum/architecture/service/post"
	"real-time-forum/architecture/service/post_comment"
	"real-time-forum/architecture/service/post_comment_vote"
	"real-time-forum/architecture/service/post_vote"
	"real-time-forum/architecture/service/session"
	"real-time-forum/architecture/service/user"
)

type Service struct {
	User            models.UserService
	Post            models.PostService
	PostVote        models.PostVoteService
	Category        models.CategoryService
	PostComment     models.PostCommentService
	PostCommentVote models.PostCommentVoteService
	Session         models.SessionService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:            user.NewUserService(repo.User),
		Post:            post.NewPostService(repo.Post),
		PostVote:        post_vote.NewPostVoteService(repo.PostVote),
		Category:        category.NewPostCategoryService(repo.Category),
		PostComment:     post_comment.NewPostCommentService(repo.PostComment),
		PostCommentVote: post_comment_vote.NewPostCommentVoteService(repo.PostCommentVote),
		Session:         session.NewSessionService(repo.Session),
	}
}
