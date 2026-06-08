package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           int64     `json:"-"`
	PublicID     uuid.UUID `json:"public_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Firstname    string    `json:"first_name"`
	Lastname     string    `json:"last_name"`
	PhoneNumber  string    `json:"phone_number"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"-"`
}

// FOR DATABASE INTERACTION LOGIC
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByPublicID(ctx context.Context, publicID uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, publicID uuid.UUID) error
}

// FOR BUSINESS LOGIC
type UserService interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email, password string) (*User, error)
	GetByPublicID(ctx context.Context, publicID uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, publicID uuid.UUID) error
}

// FOR REGISTRATION PART
type RegisterUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required, min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// FOR LOGIN PART
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required, min=8"`
}

// FOR RESPONSE PART (USER PROFILE IN THE APP. LIKE PROFILE PAGE)
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Status      string    `json:"status"`
}
