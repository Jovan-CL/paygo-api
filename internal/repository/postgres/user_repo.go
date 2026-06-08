package postgres

import (
	"context"
	"database/sql"
	"paygo-api/internal/domain"

	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *domain.User) error {

	query := ` INSERT INTO users (email, password_hash, first_name, last_name, phone_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, public_id, status, created_at, updated_at`

	err := r.DB.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Firstname, user.Lastname, user.PhoneNumber).Scan(&user.ID, &user.PublicID, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User

	query := `SELECT id, public_id, email, password_hash, first_name, last_name, phone_number, status, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL`

	err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.PublicID, &user.Email, &user.PasswordHash, &user.Firstname, &user.Lastname, &user.PhoneNumber, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindUserByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.User, error) {

	var user domain.User

	query := `SELECT id, public_id, email, password_hash, first_name, last_name, phone_number, status, created_at, updated_at FROM users WHERE public_id = $1 AND deleted_at IS NULL`

	err := r.DB.QueryRowContext(ctx, query, publicID).Scan(&user.ID, &user.PublicID, &user.Email, &user.PasswordHash, &user.Firstname, &user.Lastname, &user.PhoneNumber, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}

	return &user, nil
}
