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

	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, public_id, status, created_at, updated_at`

	err := r.DB.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Firstname, user.Lastname, user.PhoneNumber).Scan(&user.ID, &user.PublicID, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// 1. Create an empty instance of the user struct to hold the database data
	var user domain.User

	// 2. Select the critical fields your Service Layer needs for login authentication
	query := `
		SELECT id, public_id, email, password_hash, first_name, last_name, phone_number, status 
		FROM users 
		WHERE email = $1 
		AND deleted_at IS NULL`

	// 3. Match your SQL fields sequentially directly to your Go struct pointers
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.PublicID,
		&user.Email,
		&user.PasswordHash,
		&user.Firstname,
		&user.Lastname,
		&user.PhoneNumber,
		&user.Status,
	)

	// 4. Handle errors cleanly
	if err != nil {
		// If no row was found, return nil for the user and pass the standard sql error up
		return nil, err
	}

	// 5. Return the memory pointer to the populated user struct along with no error
	return &user, nil
}

func (r *PostgresUserRepository) FindUserByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.User, error) {

	var user domain.User

	query := `
		SELECT id, public_id, email, password_hash, first_name, last_name, phone_number, status, created_at, updated_at 
		FROM users 
		WHERE public_id = $1 
		AND deleted_at IS NULL`

	err := r.DB.QueryRowContext(ctx, query, publicID).Scan(&user.ID, &user.PublicID, &user.Email, &user.PasswordHash, &user.Firstname, &user.Lastname, &user.PhoneNumber, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows || err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
			UPDATE users SET first_name = $1, last_name = $2, phone_number = $3, updated_at = NOW() 
			FROM users 
			WHERE public_id = $4 
			AND deleted_at IS NULL`

	_, err := r.DB.ExecContext(ctx, query, user.Firstname, user.Lastname, user.PhoneNumber, user.PublicID)
	return err
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, publicID uuid.UUID) error {
	query := `
			UPDATE deleted_at = NOW(), status = 'deactivated' 
			WHERE public_id = $1 
			AND deleted_at IS NULL`

	_, err := r.DB.ExecContext(ctx, query, publicID)

	return err
}
