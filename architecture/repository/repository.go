package repository

import (
	"database/sql"
	"real-time-forum/architecture/models"
	"real-time-forum/architecture/repository/category"
	"real-time-forum/architecture/repository/post"
	"real-time-forum/architecture/repository/post_comment"
	"real-time-forum/architecture/repository/post_comment_vote"
	"real-time-forum/architecture/repository/post_vote"
	"real-time-forum/architecture/repository/private_message"
	"real-time-forum/architecture/repository/session"
	"real-time-forum/architecture/repository/user"
)

type Repository struct {
	User            models.UserRepo
	Post            models.PostRepo
	PostVote        models.PostVoteRepo
	Category        models.CategoryRepo
	PostComment     models.PostCommentRepo
	PostCommentVote models.PostCommentVoteRepo
	PrivateMessage  models.PrivateMessageRepo
	Session         models.SessionRepo
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		User:            user.NewUserRepo(db),
		Post:            post.NewPostRepo(db),
		PostVote:        post_vote.NewPostVoteRepo(db),
		Category:        category.NewCategoryRepo(db),
		PostComment:     post_comment.NewPostCommentRepo(db),
		PostCommentVote: post_comment_vote.NewPostCommentVoteRepo(db),
		PrivateMessage:  private_message.NewPrivateMessageRepo(db),
		Session:         session.NewSessionRepo(db),
	}
}
