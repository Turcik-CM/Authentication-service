package service

import (
	pb "auth-service/genproto/user"
	"auth-service/pkg/hashing"
	"auth-service/storage"
	"context"
	"log/slog"
)

type UserServices interface {
	Create(ctx context.Context, in *pb.CreateRequest) (*pb.UserResponse, error)
	GetProfile(ctx context.Context, in *pb.Id) (*pb.GetProfileResponse, error)
	UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UserResponse, error)
	ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)
	ChangeProfileImage(ctx context.Context, in *pb.URL) (*pb.Void, error)
	FetchUsers(ctx context.Context, in *pb.Filter) (*pb.UserResponses, error)
	PostAdd(ctx context.Context, in *pb.Id) (*pb.Void, error) // MASHI NARSA YOQ
	ListOfFollowing(ctx context.Context, in *pb.Id) (*pb.Followings, error)
	ListOfFollowers(ctx context.Context, in *pb.Id) (*pb.Followers, error)
	ListOfFollowingByUsername(ctx context.Context, in *pb.Id) (*pb.Followings, error)
	ListOfFollowersByUsername(ctx context.Context, in *pb.Id) (*pb.Followers, error)
	PostDelete(ctx context.Context, in *pb.Id) (*pb.Void, error) //MASHIYAM YOQ
	DeleteUser(ctx context.Context, in *pb.Id) (*pb.Void, error)

	Follow(ctx context.Context, in *pb.FollowReq) (*pb.FollowRes, error)
	Unfollow(ctx context.Context, in *pb.FollowReq) (*pb.DFollowRes, error)
	GetUserFollowers(ctx context.Context, in *pb.Id) (*pb.Count, error)
	GetUserFollows(ctx context.Context, in *pb.Id) (*pb.Count, error)
	MostPopularUser(ctx context.Context, in *pb.Void) (*pb.UserResponse, error)
}

type UserService struct {
	pb.UnimplementedUserServiceServer
	st  storage.UserStorage
	log *slog.Logger
}

func NewUserService(st storage.UserStorage, logger *slog.Logger) *UserService {
	return &UserService{
		st:  st,
		log: logger,
	}
}

func (us *UserService) Create(ctx context.Context, in *pb.CreateRequest) (*pb.UserResponse, error) {
	hash, err := hashing.HashPassword(in.Password)
	if err != nil {
		us.log.Error("failed to hash password", "error", err)
		return nil, err
	}

	in.Password = hash

	res, err := us.st.Create(in)
	if err != nil {
		us.log.Error("failed to create user", "error", err)
		return nil, err
	}

	return res, nil
}

func (us *UserService) GetProfile(ctx context.Context, in *pb.Id) (*pb.GetProfileResponse, error) {
	res, err := us.st.GetProfile(in)
	if err != nil {
		us.log.Error("failed to get user", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	res, err := us.st.UpdateProfile(in)
	if err != nil {
		us.log.Error("failed to update user", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	hash, err := hashing.HashPassword(in.NewPassword)
	if err != nil {
		us.log.Error("Failed to hash password", "error", err)
		return nil, err
	}

	in.NewPassword = hash

	res, err := us.st.ChangePassword(in)
	if err != nil {
		us.log.Error("failed to change password", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) ChangeProfileImage(ctx context.Context, in *pb.URL) (*pb.Void, error) {
	res, err := us.st.ChangeProfileImage(in)
	if err != nil {
		us.log.Error("failed to change profile image", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) FetchUsers(ctx context.Context, in *pb.Filter) (*pb.UserResponses, error) {
	res, err := us.st.FetchUsers(in)
	if err != nil {
		us.log.Error("failed to fetch users", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) ListOfFollowing(ctx context.Context, in *pb.Id) (*pb.Followings, error) {
	res, err := us.st.ListOfFollowing(in)
	if err != nil {
		us.log.Error("failed to list following", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) ListOfFollowers(ctx context.Context, in *pb.Id) (*pb.Followers, error) {
	res, err := us.st.ListOfFollowers(in)
	if err != nil {
		us.log.Error("failed to list followers", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) ListOfFollowingByUsername(ctx context.Context, in *pb.Id) (*pb.Followings, error) {
	res, err := us.st.ListOfFollowingByUsername(in)
	if err != nil {
		us.log.Error("failed to list following", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) ListOfFollowersByUsername(ctx context.Context, in *pb.Id) (*pb.Followers, error) {
	res, err := us.st.ListOfFollowersByUsername(in)
	if err != nil {
		us.log.Error("failed to list followers", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) DeleteUser(ctx context.Context, in *pb.Id) (*pb.Void, error) {
	res, err := us.st.DeleteUser(in)
	if err != nil {
		us.log.Error("failed to delete user", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) Follow(ctx context.Context, in *pb.FollowReq) (*pb.FollowRes, error) {
	res, err := us.st.Follow(in)
	if err != nil {
		us.log.Error("failed to follow", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) Unfollow(ctx context.Context, in *pb.FollowReq) (*pb.DFollowRes, error) {
	res, err := us.st.Unfollow(in)
	if err != nil {
		us.log.Error("failed to unfollow", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) GetUserFollowers(ctx context.Context, in *pb.Id) (*pb.Count, error) {
	res, err := us.st.GetUserFollowers(in)
	if err != nil {
		us.log.Error("failed to get user followers", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *UserService) GetUserFollows(ctx context.Context, in *pb.Id) (*pb.Count, error) {
	res, err := us.st.GetUserFollows(in)
	if err != nil {
		us.log.Error("failed to get user follows", "error", err)
		return nil, err
	}
	return res, nil
}
func (us *UserService) MostPopularUser(ctx context.Context, in *pb.Void) (*pb.UserResponse, error) {
	res, err := us.st.MostPopularUser(in)
	if err != nil {
		us.log.Error("failed to most popular user", "error", err)
		return nil, err
	}
	return res, nil
}
