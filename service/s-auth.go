package service

import (
	pb "auth-service/genproto/user"
	"context"
)

func (s *AuthService) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res, err := s.st.Register(in)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *AuthService) LoginEmail(ctx context.Context, in *pb.LoginEmailRequest) (*pb.LoginResponse, error) {
	res, err := s.st.LoginEmail(in)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *AuthService) LoginUsername(ctx context.Context, in *pb.LoginUsernameRequest) (*pb.LoginResponse, error) {
	res, err := s.st.LoginUsername(in)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, in *pb.Email) (*pb.GetProfileResponse, error) {
	res, err := s.st.GetUserByEmail(in)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *AuthService) RegisterAdmin(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	res, err := s.st.RegisterAdmin(in)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, err
	}
	return res, nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, in *pb.UpdatePasswordReq) (*pb.Message, error) {
	res, err := s.st.UpdatePassword(in)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, err
	}
	return res, nil
}
