package private_message

import (
	"database/sql"
	"real-time-forum/architecture/models"
)

type PrivateMessageRepo struct {
	db *sql.DB
}

func NewPrivateMessageRepo(db *sql.DB) models.PrivateMessageRepo {
	return &PrivateMessageRepo{db: db}
}
