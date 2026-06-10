package service

import (
	"context"
	"database/sql"
	"errors"
	"paygo-api/internal/domain"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: repo, // Injects your Postgres tool right into the brain
	}
}

func (s *userService) Register(ctx context.Context, user *domain.User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	passHashed, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to process account security tokens.")
	}
	user.PasswordHash = string(passHashed)

	user.Status = "active"

	return s.userRepo.CreateUser(ctx, user)
}

func (s *userService) Login(ctx context.Context, email string, rawPassword string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	dbUser, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("Invalid Email or password!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(rawPassword))
	if err != nil {
		return nil, errors.New("Invalid Email or password!")
	}
	return dbUser, nil
}

func (s *userService) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.User, error) {

	dbUser, err := s.userRepo.FindUserByPublicID(ctx, publicID)
	if err != nil {
		return nil, errors.New("Failed to retrieve user data")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("User profile not found.")
	}
	if dbUser.Status == "suspended" {
		return nil, errors.New("This user account is suspended due to compliance violations")
	}

	if dbUser.Status == "deactivated" {
		softDeleteDate := "unkown time"
		if !dbUser.DeletedAt.IsZero() {
			softDeleteDate = dbUser.DeletedAt.Format("Jan 02, 2006 - 03:04 PM")
		}

		return nil, errors.New("this user account was deactivated at " + softDeleteDate)
	}

	dbUser.Firstname = strings.TrimSpace(dbUser.Firstname)
	dbUser.Lastname = strings.TrimSpace(dbUser.Lastname)

	return dbUser, nil
}

func (s *userService) Delete(ctx context.Context, publicID uuid.UUID) error {
	// Pass the delete signal directly down to your database repository worker
	return s.userRepo.DeleteUser(ctx, publicID)
}

func (s *userService) Update(ctx context.Context, user *domain.User) error {
	return s.userRepo.UpdateUser(ctx, user)
}
