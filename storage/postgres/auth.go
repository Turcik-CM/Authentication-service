package postgres

import (
	pb "auth-service/genproto/user"
	"auth-service/pkg/hashing"
	"auth-service/pkg/token"
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

func (a *AuthRepo) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var id string
	var flag string
	query := `INSERT INTO users (phone, email, password, first_name, nationality, last_name, username, bio) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			  RETURNING id`
	err := a.db.QueryRow(query, in.Phone, in.Email, in.Password, in.FirstName, in.Nationality, in.LastName, in.Username, in.Bio).Scan(&id)
	if err != nil {
		return &pb.RegisterResponse{}, err
	}

	return &pb.RegisterResponse{
		Id:    id,
		Email: in.Email,
		Flag:  flag,
	}, nil
}

func (a *AuthRepo) GetUserByEmail(in *pb.Email) (*pb.GetProfileResponse, error) {
	query := `SELECT id, created_at FROM users WHERE email = $1 AND deleted_at=0`

	var user pb.GetProfileResponse
	err := a.db.QueryRowContext(context.Background(), query, in.Email).Scan(&user.UserId, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (a *AuthRepo) LoginEmail(in *pb.LoginEmailRequest) (*pb.LoginResponse, error) {
	res := pb.LoginResponse1{}
	query := `SELECT id, email, password, role, username FROM users WHERE email = $1 AND deleted_at = 0`
	err := a.db.Get(&res, query, in.Email)
	if err != nil {
		return &pb.LoginResponse{}, err
	}

	check := hashing.CheckPasswordHash(res.Password, in.Password)
	if !check {
		return &pb.LoginResponse{}, errors.New("invalid password")
	}

	reqToken := pb.LoginResponse1{
		Id:       res.Id,
		Email:    res.Email,
		UserName: res.UserName,
		Role:     res.Role,
		Country:  res.Country,
	}

	accessToken, err := token.GenerateAccessToken(&reqToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken(&reqToken)
	if err != nil {
		return nil, err
	}

	www := pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       res.Id,
	}

	return &www, nil
}

func (a *AuthRepo) LoginUsername(in *pb.LoginUsernameRequest) (*pb.LoginResponse, error) {
	res := pb.LoginResponse1{}
	query := `SELECT id, email, password, role, username FROM users WHERE username = $1 AND deleted_at = 0`
	err := a.db.Get(&res, query, in.Username)
	if err != nil {
		return &pb.LoginResponse{}, err
	}

	check := hashing.CheckPasswordHash(res.Password, in.Password)
	if !check {
		return &pb.LoginResponse{}, errors.New("invalid password")
	}

	reqToken := pb.LoginResponse1{
		Id:       res.Id,
		Email:    res.Email,
		UserName: res.UserName,
		Role:     res.Role,
		Country:  res.Country,
	}

	accessToken, err := token.GenerateAccessToken(&reqToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken(&reqToken)
	if err != nil {
		return nil, err
	}

	www := pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       res.Id,
	}

	return &www, nil
}

func (a *AuthRepo) RegisterAdmin(in *pb.Message) (*pb.Message, error) {
	query := `INSERT INTO users (email, password, role, nationality,first_name,last_name,username,bio,phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	var id string
	err := a.db.QueryRow(query, "admiN", in.Message, "c-admin", "Uzbek", "adminchikov", "admin", "admin", "admin", "admin").Scan(&id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *AuthRepo) UpdatePassword(in *pb.UpdatePasswordReq) (*pb.Message, error) {
	query := `UPDATE users SET password = $1 WHERE id = $2 AND deleted_at = 0`

	result, err := a.db.ExecContext(context.Background(), query, in.Password, in.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return nil, nil
}
