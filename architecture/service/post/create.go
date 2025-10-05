package post

import (
	"errors"
	"fmt"
	"real-time-forum/architecture/models"
	"time"
)

func (p *PostService) Create(post *models.Post) (int64, error) {
	Prepare(post)

	if ValidateTitle(post) != nil {
		return -1, ErrInvalidTitleLength
	} else if ValidateContent(post) != nil {
		return -1, ErrInvalidContentLength
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt

	postId, err := p.repo.Create(post)
	switch {
	case err == nil:
	case errors.Is(err, ErrInvalidTitleLength):
		return -1, ErrInvalidTitleLength
	case err != nil:
		return -1, fmt.Errorf("p.repo.Create: %w", err)
	}
	return postId, nil

}
