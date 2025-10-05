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
	if err := ValidateNickname(user); err != nil {
		return ErrInvalidNickname
	}
	if err := ValidateEmail(user); err != nil {
		return ErrInvalidEmail
	}
	if err := ValidateAge(user); err != nil {
		return ErrInvalidAge
	}
	if err := ValidateGender(user); err != nil {
		return ErrInvalidGender
	}

	// Хешируем пароль если он изменился
	if user.Password != "" {
		if err := HashPassword(user); err != nil {
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
