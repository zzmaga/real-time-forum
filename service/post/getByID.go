package post

import (
	"errors"
	"fmt"
	"real-time-forum/architecture/models"
)

func (p *PostService) GetByID(id int64) (*models.Post, error) {
	post, err := p.repo.GetByID(id)
	switch {
	case err == nil:
		return post, nil
	case errors.Is(err, ErrNotFound):
		return nil, ErrNotFound
	}
	return nil, fmt.Errorf("p.repo.GetByID: %w", err)
}
