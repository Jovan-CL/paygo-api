package service

import (
	"context"
	"errors"
	"paygo-api/internal/domain"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo domain.UserRepository
}

func (s *userService) Register(ctx context.Context, user domain.User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	passHashed, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to process account security tokens.")
	}
	user.PasswordHash = string(passHashed)

	user.Status = "active"

	return s.userRepo.CreateUser(ctx, &user)
}
