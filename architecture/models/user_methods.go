package models

import (
	"fmt"
	"net/mail"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func (u *User) ValidateNickname() error {
	if lng := len([]rune(u.Nickname)); lng < 1 || 32 < lng {
		return fmt.Errorf("nickname: invalid length (%d)", lng)
	}
	for _, c := range u.Nickname {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_') {
			return fmt.Errorf("nickname: invalid character '%c'", c)
		}
	}
	return nil
}

func (u *User) ValidateEmail() error {
	if lng := len([]rune(u.Email)); lng < 1 || 320 < lng {
		return fmt.Errorf("email: invalid length (%d)", lng)
	}
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ValidatePassword() error {
	if len([]rune(u.Password)) < 8 {
		return fmt.Errorf("password: too short")
	}

	var hasLetter, hasDigit bool
	for _, c := range u.Password {
		if unicode.IsLetter(c) {
			hasLetter = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return fmt.Errorf("password: must contain at least one letter and one digit")
	}
	return nil
}

func (u *User) ValidateAge() error {
	if u.Age < 0 || u.Age > 150 {
		return fmt.Errorf("age: invalid value (%d), must be between 0 and 150", u.Age)
	}
	return nil
}

func (u *User) ValidateGender() error {
	validGenders := []string{"male", "female", "other"}
	for _, gender := range validGenders {
		if u.Gender == gender {
			return nil
		}
	}
	return fmt.Errorf("gender: invalid value (%s), must be one of: male, female, other", u.Gender)
}

func (u *User) HashPassword() error {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}
	u.Password = string(pass)
	return nil
}

func (u *User) CompareHashAndPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	switch err {
	case bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
	}
}
