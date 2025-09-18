package category

import (
	"real-time-forum/architecture/models"
)

type CategoryService struct {
	repo models.CategoryRepo
}

func NewPostCategoryService(repo models.CategoryRepo) *CategoryService {
	return &CategoryService{repo}
}
