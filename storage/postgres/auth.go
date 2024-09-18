package postgres

import (
	"auth-service/pkg/models"
	"auth-service/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) storage.AuthStorage {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) Register(in models.RegisterRequest) (models.RegisterResponse, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return models.RegisterResponse{}, err
	}

	var id string
	query := `INSERT INTO users (phone, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = a.db.QueryRow(query, in.Phone, in.Email, in.Password).Scan(&id)
	if err != nil {
		return models.RegisterResponse{}, err
	}

	query1 := `INSERT INTO user_profile (user_id, first_name, last_name, username, nationality, bio) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = a.db.Query(query1, id, in.FirstName, in.LastName, in.Username, in.Nationality, in.Bio)
	if err != nil {
		return models.RegisterResponse{}, err
	}

	err = tx.Commit()
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
func (a *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error) {
	query := `SELECT id ,created_at FROM users WHERE email = $1 AND deleted_at=0`

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
func (a *AuthRepo) LoginEmail(in models.LoginEmailRequest) (models.LoginResponse, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return models.LoginResponse{}, err
	}

	res := models.LoginResponse{}

	query := `SELECT id, email, password FROM users WHERE email = $1 and deleted_at = 0`
	err = a.db.Get(&res, query, in.Email)
	if err != nil {
		return models.LoginResponse{}, err
	}

	query1 := `SELECT role, username FROM user_profile WHERE user_id = $1`
	err = a.db.Get(&res, query1, res.Id)

	err = tx.Commit()
	if err != nil {
		return models.LoginResponse{}, err
	}

	return res, nil
}
func (a *AuthRepo) LoginUsername(in models.LoginUsernameRequest) (models.LoginResponse, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return models.LoginResponse{}, err
	}

	var id string
	var role string
	query := `SELECT role, user_id FROM user_profile WHERE username = $1`
	err = a.db.QueryRow(query, in.Username).Scan(&role, &id)
	if err != nil {
		return models.LoginResponse{}, err
	}

	res := models.LoginResponse{}
	query1 := `SELECT password, email FROM users WHERE id = $1 and deleted_at = 0`
	err = a.db.Get(&res, query1, id)
	if err != nil {
		return models.LoginResponse{}, err
	}
	log.Println()
	err = tx.Commit()
	if err != nil {
		return models.LoginResponse{}, err
	}

	return models.LoginResponse{
		Id:       id,
		Email:    res.Email,
		Username: in.Username,
		Password: res.Password,
		Role:     role,
	}, nil
}
func (a *AuthRepo) RegisterAdmin(ctx context.Context, pass string) error {
	tx, err := a.db.Begin()
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id`
	var id string
	err = a.db.QueryRow(query, "admiN", pass).Scan(&id)
	if err != nil {
		err = tx.Rollback()
		return err
	}

	query = `INSERT INTO user_profile (user_id,role)
	VALUES ($1, $2) 
`
	err = a.db.QueryRow(query, id, "c-admin").Err()
	if err != nil {
		err = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (a *AuthRepo) UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error {
	query := `update users set password=$1 where id=$2 and deleted_at=0`

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
