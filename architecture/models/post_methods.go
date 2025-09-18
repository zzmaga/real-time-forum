package models

import (
	"fmt"
	"strings"
)

func (p *Post) ValidateTitle() error {
	if lng := len([]rune(p.Title)); lng < 1 || lng > 100 {
		return fmt.Errorf("title: invalid lenght (%d)", lng)
	}
	return nil
}

func (p *Post) ValidateContent() error {
	if lng := len([]rune(p.Content)); lng < 1 || lng > 3000 {
		return fmt.Errorf("content: invalid lenght (%d)", lng)
	}
	return nil
}

func (p *Post) PrepareTitle() {
	p.Title = strings.Trim(p.Title, " ")
}

func (p *Post) PrepareContent() {
	p.Content = strings.Trim(p.Content, " ")
}

func (p *Post) Prepare() {
	p.PrepareTitle()
	p.PrepareContent()
}
