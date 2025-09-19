package user

import (
	"errors"
	"fmt"
	"real-time-forum/architecture/models"
	ruser "real-time-forum/architecture/repository/user"
	"time"
)

func (u *UserService) Update(user *models.User) error {
	// Валидация полей
	if err := user.ValidateNickname(); err != nil {
		return ErrInvalidNickname
	}
	if err := user.ValidateEmail(); err != nil {
		return ErrInvalidEmail
	}
	if err := user.ValidateAge(); err != nil {
		return ErrInvalidAge
	}
	if err := user.ValidateGender(); err != nil {
		return ErrInvalidGender
	}

	// Хешируем пароль если он изменился
	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return fmt.Errorf("user.HashPassword: %w", err)
		}
	}

	// Устанавливаем время обновления
	user.UpdatedAt = time.Now()

	// Обновляем в репозитории
	err := u.repo.Update(user)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ruser.ErrNotFound):
		return ErrNotFound
	}
	return fmt.Errorf("u.repo.Update: %w", err)
}
