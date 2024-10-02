package service

import (
	pb "auth-service/genproto/user"
	"auth-service/pkg/hashing"
	"context"
	"log"
)

func (us *AuthService) Create(ctx context.Context, in *pb.CreateRequest) (*pb.UserResponse, error) {
	hash, err := hashing.HashPassword(in.Password)
	if err != nil {
		us.log.Error("failed to hash password", "error", err)
		return nil, err
	}

	in.Password = hash

	res, err := us.user.Create(in)
	if err != nil {
		us.log.Error("failed to create user", "error", err)
		return nil, err
	}

	return res, nil
}

func (us *AuthService) GetProfile(ctx context.Context, in *pb.Id) (*pb.GetProfileResponse, error) {
	res, err := us.user.GetProfile(in)
	if err != nil {
		us.log.Error("failed to get user", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	res, err := us.user.UpdateProfile(in)
	if err != nil {
		us.log.Error("failed to update user", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	hash, err := hashing.HashPassword(in.NewPassword)
	if err != nil {
		us.log.Error("Failed to hash password", "error", err)
		return nil, err
	}

	in.NewPassword = hash

	res, err := us.user.ChangePassword(in)
	if err != nil {
		us.log.Error("failed to change password", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) ChangeProfileImage(ctx context.Context, in *pb.URL) (*pb.Void, error) {
	res, err := us.user.ChangeProfileImage(in)
	if err != nil {
		us.log.Error("failed to change profile image", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) FetchUsers(ctx context.Context, in *pb.Filter) (*pb.UserResponses, error) {
	res, err := us.user.FetchUsers(in)
	if err != nil {
		us.log.Error("failed to fetch users", "error", err)
		return nil, err
	}
	log.Println("hello")
	return res, nil
}

func (us *AuthService) ListOfFollowing(ctx context.Context, in *pb.Id) (*pb.Follows, error) {
	res, err := us.user.ListOfFollowing(in)
	if err != nil {
		us.log.Error("failed to list following", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) ListOfFollowers(ctx context.Context, in *pb.Id) (*pb.Follows, error) {
	res, err := us.user.ListOfFollowers(in)
	if err != nil {
		us.log.Error("failed to list followers", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) DeleteUser(ctx context.Context, in *pb.Id) (*pb.Void, error) {
	res, err := us.user.DeleteUser(in)
	if err != nil {
		us.log.Error("failed to delete user", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) Follow(ctx context.Context, in *pb.FollowReq) (*pb.FollowRes, error) {
	res, err := us.user.Follow(in)
	if err != nil {
		us.log.Error("failed to follow", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) Unfollow(ctx context.Context, in *pb.FollowReq) (*pb.DFollowRes, error) {
	res, err := us.user.Unfollow(in)
	if err != nil {
		us.log.Error("failed to unfollow", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) GetUserFollowers(ctx context.Context, in *pb.Id) (*pb.Count, error) {
	res, err := us.user.GetUserFollowers(in)
	if err != nil {
		us.log.Error("failed to get user followers", "error", err)
		return nil, err
	}
	return res, nil
}

func (us *AuthService) GetUserFollows(ctx context.Context, in *pb.Id) (*pb.Count, error) {
	res, err := us.user.GetUserFollows(in)
	if err != nil {
		us.log.Error("failed to get user follows", "error", err)
		return nil, err
	}
	return res, nil
}
func (us *AuthService) MostPopularUser(ctx context.Context, in *pb.Void) (*pb.UserResponse, error) {
	res, err := us.user.MostPopularUser(in)
	if err != nil {
		us.log.Error("failed to most popular user", "error", err)
		return nil, err
	}
	return res, nil
}
