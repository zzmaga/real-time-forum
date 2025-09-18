package post

import "database/sql"

type PostRepo struct {
	DB *sql.DB
}

// Constructor
func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{DB: db}
}
