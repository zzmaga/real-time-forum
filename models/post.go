package models

import "time"

type Post struct {
	Id        int64
	Title     string
	Content   string
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time

	WUser       *User
	WUserVote   int8  // -1 0 1
	WVoteUp     int64 // Like
	WVoteDown   int64 // Dislike
	WCategories []*Category
	WComments   []*PostComment
}

type PostService interface {
	Create(post *Post) (int64, error)
	Update(post *Post) error
	GetAll(offset, limit int64) ([]*Post, error)
	GetByID(id int64) (*Post, error)
	GetByIDs(ids []int64) ([]*Post, error)
	GetByUserID(userId, offset, limit int64) ([]*Post, error)
	DeleteByID(id int64) error
}

type PostRepo interface {
	Create(post *Post) (int64, error)
	Update(post *Post) error
	GetAll(offset, limit int64) ([]*Post, error)
	GetByID(id int64) (*Post, error)
	GetByIDs(ids []int64) ([]*Post, error)
	GetByUserID(userId, offset, limit int64) ([]*Post, error)
	DeleteByID(id int64) error
}
