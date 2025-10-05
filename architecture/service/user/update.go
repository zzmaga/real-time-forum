package user

import (
	"errors"
	"fmt"
	"real-time-forum/architecture/models"
	"time"
)

func (u *UserService) Update(user *models.User) error {
	// Валидация
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

	// Хешируем пароль
	if user.Password != "" {
		if err := HashPassword(user); err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
	}

	user.UpdatedAt = time.Now()

	// Обновляем
	err := u.repo.Update(user)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}

	// Логируем/пробрасываем
	return fmt.Errorf("repo update: %w", err)
}
