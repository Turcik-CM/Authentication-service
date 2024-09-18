package service

import (
	"auth-service/pkg/config"
	"auth-service/pkg/hashing"
	"auth-service/pkg/models"
	"auth-service/pkg/token"
	"auth-service/storage"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
)

type AuthService interface {
	Register(in models.RegisterRequest) (models.RegisterResponse, error)
	LoginEmail(in models.LoginEmailRequest) (models.Tokens, error)
	LoginUsername(in models.LoginUsernameRequest) (models.Tokens, error)
	GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error)
	RegisterAdmin(ctx context.Context) error
	UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error
}

func NewAuthService(st storage.AuthStorage, logger *slog.Logger) AuthService {
	return &authService{st, logger}
}

type authService struct {
	st  storage.AuthStorage
	log *slog.Logger
}

func (a *authService) Register(in models.RegisterRequest) (models.RegisterResponse, error) {
	hash, err := hashing.HashPassword(in.Password)
	if err != nil {
		a.log.Error("Failed to hash password", "error", err)
		return models.RegisterResponse{}, err
	}

	in.Password = hash

	res, err := a.st.Register(in)
	if err != nil {
		a.log.Error("Failed to register user", "error", err)
		return models.RegisterResponse{}, err
	}

	return res, nil
}

func (a *authService) LoginEmail(in models.LoginEmailRequest) (models.Tokens, error) {
	res, err := a.st.LoginEmail(in)
	if err != nil {
		a.log.Error("Failed to login", "error", err)
		return models.Tokens{}, err
	}

	check := hashing.CheckPasswordHash(res.Password, in.Password)
	if !check {
		a.log.Error("Invalid password")
		return models.Tokens{}, errors.New("Invalid password")
	}

	refreshToken, err := token.GenerateRefreshToken(res)
	if err != nil {
		a.log.Error("Failed to generate refresh token", "error", err)
		return models.Tokens{}, err
	}

	accessToken, err := token.GenerateAccessToken(res)
	if err != nil {
		a.log.Error("Failed to generate access token", "error", err)
		return models.Tokens{}, err
	}

	response := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (a *authService) LoginUsername(in models.LoginUsernameRequest) (models.Tokens, error) {
	res, err := a.st.LoginUsername(in)
	if err != nil {
		a.log.Error("Failed to login", "error", err)
		return models.Tokens{}, err
	}

	check := hashing.CheckPasswordHash(res.Password, in.Password)
	log.Println(check)
	if !check {
		a.log.Error("Invalid password")
		log.Println("\n----------", in.Password, res.Password, "\n---------")
		return models.Tokens{}, errors.New("Invalid password")
	}

	refreshToken, err := token.GenerateRefreshToken(res)
	if err != nil {
		a.log.Error("Failed to generate refresh token", "error", err)
		return models.Tokens{}, err
	}

	accessToken, err := token.GenerateAccessToken(res)
	if err != nil {
		a.log.Error("Failed to generate access token", "error", err)
		return models.Tokens{}, err
	}

	response := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}
func (a *authService) GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error) {
	a.log.Info("Getting user user by email")
	res, err := a.st.GetUserByEmail(ctx, email)
	if err != nil {
		a.log.Error(err.Error())
		return nil, err
	}
	return res, nil
}

func (a *authService) RegisterAdmin(ctx context.Context) error {

	hash, err := hashing.HashPassword(config.Load().ADMIN_PASSWORD)
	if err != nil {
		a.log.Error("Failed to hash password", "error", err)
		return err
	}

	err = a.st.RegisterAdmin(ctx, hash)
	if err != nil {
		a.log.Error("Failed to register admin", "error", err)
		return err
	}

	return nil
}

func (a *authService) UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error {
	hash, err := hashing.HashPassword(req.Password)
	if err != nil {
		a.log.Error("Failed to hash password", "error", err)
		return err
	}

	req.Password = hash

	err = a.st.UpdatePassword(ctx, req)
	if err != nil {
		a.log.Error(fmt.Sprintf("Error update pasword: %v", err))
		return err
	}
	return nil
}
