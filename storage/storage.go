package storage

import (
	pb "auth-service/genproto/user"
)

type AuthStorage interface {
	Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error)
	LoginEmail(in *pb.LoginEmailRequest) (*pb.LoginResponse, error)
	LoginUsername(in *pb.LoginUsernameRequest) (*pb.LoginResponse, error)
	GetUserByEmail(in *pb.Email) (*pb.GetProfileResponse, error)
	RegisterAdmin(in *pb.Message) (*pb.Message, error)
	UpdatePassword(in *pb.UpdatePasswordReq) (*pb.Message, error)
}

type UserStorage interface {
	Create(in *pb.CreateRequest) (*pb.UserResponse, error)
	GetProfile(in *pb.Id) (*pb.GetProfileResponse, error)
	UpdateProfile(in *pb.UpdateProfileRequest) (*pb.UserResponse, error)
	ChangePassword(in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)
	ChangeProfileImage(in *pb.URL) (*pb.Void, error)
	FetchUsers(in *pb.Filter) (*pb.UserResponses, error)
	ListOfFollowing(in *pb.Id) (*pb.Follows, error)
	ListOfFollowers(in *pb.Id) (*pb.Follows, error)
	DeleteUser(in *pb.Id) (*pb.Void, error)

	Follow(in *pb.FollowReq) (*pb.FollowRes, error)
	Unfollow(in *pb.FollowReq) (*pb.DFollowRes, error)
	GetUserFollowers(in *pb.Id) (*pb.Count, error)
	GetUserFollows(in *pb.Id) (*pb.Count, error)
	MostPopularUser(in *pb.Void) (*pb.UserResponse, error)

	AddNationality(in *pb.Nationality) (*pb.Void, error)
	GetNationalityById(in *pb.NId) (*pb.Nationality, error)
	ListNationalities(in *pb.Pagination) (*pb.Nationalities, error)
	UpdateNationality(in *pb.Nationality) (*pb.Void, error)
	DeleteNationality(in *pb.NId) (*pb.Void, error)
}
