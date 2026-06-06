package domain

import (
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

type RegisterUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required, min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required, min=8"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Status      string    `json:"status"`
}
