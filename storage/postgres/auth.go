package postgres

import (
	"auth-service/pkg/models"
	"auth-service/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) storage.AuthStorage {
	return &AuthRepo{
		db: db,
	}
}

// Register a new user with data merged into the 'users' table
func (a *AuthRepo) Register(in models.RegisterRequest) (models.RegisterResponse, error) {
	var id string
	query := `INSERT INTO users (phone, email, password, first_name, last_name, username, nationality, bio) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			  RETURNING id`
	err := a.db.QueryRow(query, in.Phone, in.Email, in.Password, in.FirstName, in.LastName, in.Username, in.Nationality, in.Bio).Scan(&id)
	if err != nil {
		return models.RegisterResponse{}, err
	}

	return models.RegisterResponse{
		Id:          id,
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		Phone:       in.Phone,
		Username:    in.Username,
		Nationality: in.Nationality,
		Bio:         in.Bio,
	}, nil
}

// GetUserByEmail retrieves user data by email
func (a *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error) {
	query := `SELECT id, created_at FROM users WHERE email = $1 AND deleted_at=0`

	var user models.GetProfileResponse
	err := a.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Login using email and get user info
func (a *AuthRepo) LoginEmail(in models.LoginEmailRequest) (models.LoginResponse, error) {
	res := models.LoginResponse{}
	query := `SELECT id, email, password, role, username FROM users WHERE email = $1 AND deleted_at = 0`
	err := a.db.Get(&res, query, in.Email)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return res, nil
}

// Login using username and get user info
func (a *AuthRepo) LoginUsername(in models.LoginUsernameRequest) (models.LoginResponse, error) {
	res := models.LoginResponse{}
	query := `SELECT id, email, password, role, username FROM users WHERE username = $1 AND deleted_at = 0`
	err := a.db.Get(&res, query, in.Username)
	if err != nil {
		return models.LoginResponse{}, err
	}

	return res, nil
}

// RegisterAdmin creates a new admin user
func (a *AuthRepo) RegisterAdmin(ctx context.Context, pass string) error {
	query := `INSERT INTO users (email, password, role,first_name,last_name,username) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id string
	err := a.db.QueryRow(query, "admiN", pass, "c-admin", "adminchikov", "admin", "admin").Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword allows the user to change their password
func (a *AuthRepo) UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error {
	query := `UPDATE users SET password = $1 WHERE id = $2 AND deleted_at = 0`

	result, err := a.db.ExecContext(ctx, query, req.Password, req.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
