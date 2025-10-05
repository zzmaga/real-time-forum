package post

import (
	"fmt"
	"real-time-forum/architecture/models"
	"strings"
)

func ValidateTitle(p *models.Post) error {
	if lng := len([]rune(p.Title)); lng < 1 || lng > 100 {
		return fmt.Errorf("title: invalid lenght (%d)", lng)
	}
	return nil
}

func ValidateContent(p *models.Post) error {
	if lng := len([]rune(p.Content)); lng < 1 || lng > 3000 {
		return fmt.Errorf("content: invalid lenght (%d)", lng)
	}
	return nil
}

func PrepareTitle(p *models.Post) {
	p.Title = strings.Trim(p.Title, " ")
}

func PrepareContent(p *models.Post) {
	p.Content = strings.Trim(p.Content, " ")
}

func Prepare(p *models.Post) {
	PrepareTitle(p)
	PrepareContent(p)
}
