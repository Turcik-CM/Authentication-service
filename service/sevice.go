package service

import (
	pb "auth-service/genproto/user"
	"auth-service/storage"
	"log/slog"
)

type AuthService struct {
	st   storage.AuthStorage
	user storage.UserStorage
	log  slog.Logger
	pb.UnimplementedUserServiceServer
}

func NewAuthService(auth storage.AuthStorage, user storage.UserStorage, logger *slog.Logger) *AuthService {
	return &AuthService{st: auth, user: user, log: *logger}
}
