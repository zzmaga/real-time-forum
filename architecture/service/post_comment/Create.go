package post_comment

import (
	"fmt"
	"log"
	"real-time-forum/architecture/models"
	"strings"
	"time"
)

func (c *PostCommentService) Create(comment *models.PostComment) (int64, error) {
	PrepareContent(comment)

	if ValidateContent(comment) != nil {
		return -1, ErrInvalidContentLength
	}
	comment.CreatedAt = time.Now()
	commentId, err := c.repo.Create(comment)
	switch {
	case err == nil:
	case err != nil:
		log.Println(err.Error())
		return -1, fmt.Errorf("c.repo.Create: %w", err)
	}
	return commentId, nil
}

func ValidateContent(c *models.PostComment) error {
	if lng := len(c.Content); lng < 1 || lng > 1000 {
		return fmt.Errorf("content: invalid lenght (%d)", lng)
	}
	return nil
}

func PrepareContent(c *models.PostComment) {
	c.Content = strings.Trim(c.Content, " ")
}
