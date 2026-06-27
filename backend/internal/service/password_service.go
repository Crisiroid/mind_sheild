package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	BcryptCost int
}

func NewPasswordService(bcryptCost int) *PasswordService {
	if bcryptCost == 0 {
		bcryptCost = bcrypt.DefaultCost
	}
	return &PasswordService{
		BcryptCost: bcryptCost,
	}
}

func (s *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (s *PasswordService) CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	return nil
}
