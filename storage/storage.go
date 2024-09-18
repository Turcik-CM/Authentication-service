package storage

import (
	pb "auth-service/genproto/user"
	"auth-service/pkg/models"
	"context"
)

type AuthStorage interface {
	Register(in models.RegisterRequest) (models.RegisterResponse, error)
	LoginEmail(in models.LoginEmailRequest) (models.LoginResponse, error)
	LoginUsername(in models.LoginUsernameRequest) (models.LoginResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error)
	RegisterAdmin(ctx context.Context, pass string) error
	UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error
}

type UserStorage interface {
	Create(in *pb.CreateRequest) (*pb.UserResponse, error)
	GetProfile(in *pb.Id) (*pb.GetProfileResponse, error)
	UpdateProfile(in *pb.UpdateProfileRequest) (*pb.UserResponse, error)
	ChangePassword(in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)
	ChangeProfileImage(in *pb.URL) (*pb.Void, error)
	FetchUsers(in *pb.Filter) (*pb.UserResponses, error)
	ListOfFollowing(in *pb.Id) (*pb.Followings, error)
	ListOfFollowers(in *pb.Id) (*pb.Followers, error)
	ListOfFollowingByUsername(req *pb.Id) (*pb.Followings, error)
	ListOfFollowersByUsername(req *pb.Id) (*pb.Followers, error)
	DeleteUser(in *pb.Id) (*pb.Void, error)

	Follow(in *pb.FollowReq) (*pb.FollowRes, error)
	Unfollow(in *pb.FollowReq) (*pb.DFollowRes, error)
	GetUserFollowers(in *pb.Id) (*pb.Count, error)
	GetUserFollows(in *pb.Id) (*pb.Count, error)
	MostPopularUser(in *pb.Void) (*pb.UserResponse, error)
}
