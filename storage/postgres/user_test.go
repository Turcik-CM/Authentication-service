package postgres

import (
	pb "auth-service/genproto/user"
	"fmt"
	"github.com/jmoiron/sqlx"
	"testing"
)

func ConnectUser() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "postgres", "123321", "turk_auth")

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreate(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	rst := pb.CreateRequest{
		Email:       "hamidullox2@gmail.com",
		Password:    "hamidullox2",
		Phone:       "9997471782",
		FirstName:   "hamidullox2",
		LastName:    "hamidullox2",
		Username:    "hamidullox2",
		Bio:         "hamidullox2",
		Nationality: "hamidullox2",
	}

	user := NewUserRepo(db)

	req, err := user.Create(&rst)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	fmt.Println(req)
}

func TestGetProfile(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)
	rst := pb.Id{
		UserId: "3a27bff7-40e4-4074-a63b-b5af91211e2f",
	}

	req, err := user.GetProfile(&rst)
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}

	fmt.Println(req)
}

func TestUpdateProfile(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)
	rst := pb.UpdateProfileRequest{
		UserId:       "3a27bff7-40e4-4074-a63b-b5af91211e2f",
		FirstName:    "hamidullox4",
		LastName:     "hamidullox4",
		Bio:          "hamidullox4",
		Nationality:  "hamidullox4",
		Username:     "hamidullox4",
		ProfileImage: "hamidullox4",
		Phone:        "9997471782",
	}

	req, err := user.UpdateProfile(&rst)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}
	fmt.Println(req)
}

func TestDeleteProfile(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)
	rst := pb.Id{
		UserId: "17b04d19-ddf6-42a0-81ba-219cfd618956",
	}
	req, err := user.DeleteUser(&rst)
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}
	fmt.Println(req)
}

func TestChangePassword(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)
	rst := pb.ChangePasswordRequest{
		UserId:          "17b04d19-ddf6-42a0-81ba-219cfd618956",
		CurrentPassword: "hamidullox4",
		NewPassword:     "hamidullox5",
	}
	res, err := user.ChangePassword(&rst)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}
	fmt.Println(res)
}

func TestChangeProfileImage(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)
	rst := pb.URL{
		UserId: "17b04d19-ddf6-42a0-81ba-219cfd618956",
		Url:    "...",
	}
	req, err := user.ChangeProfileImage(&rst)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}
	fmt.Println(req)
}

func TestFetchUsers(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)
	rst := pb.Filter{
		Limit:     10,
		Page:      1,
		FirstName: "hamidullox4",
	}
	req, err := user.FetchUsers(&rst)
	if err != nil {
		t.Errorf("Failed to fetch users: %v", err)
	}
	fmt.Println(req)
}

func TestListOfFollowing(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)
	ret := pb.Id{
		UserId: "17b04d19-ddf6-42a0-81ba-219cfd618956",
	}
	req, err := user.ListOfFollowing(&ret)
	if err != nil {
		t.Errorf("Failed to list followers: %v", err)
	}
	fmt.Println(req)
}

func TestListOfFollowers(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)
	ret := pb.Id{
		UserId: "17b04d19-ddf6-42a0-81ba-219cfd618956",
	}
	req, err := user.ListOfFollowers(&ret)
	if err != nil {
		t.Errorf("Failed to list followers: %v", err)
	}
	fmt.Println(req)
}

func TestFollow(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}

	user := NewUserRepo(db)

	rst := pb.FollowReq{
		FollowingId: "3a27bff7-40e4-4074-a63b-b5af91211e2f",
		FollowerId:  "ef778e7b-059c-4117-8e8d-837a3dff0e76",
	}

	req, err := user.Follow(&rst)
	if err != nil {
		t.Errorf("Failed to follow: %v", err)
	}
	fmt.Println(req)
}

func TestUnfollow(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	rst := pb.FollowReq{
		FollowingId: "3a27bff7-40e4-4074-a63b-b5af91211e2f",
		FollowerId:  "17b04d19-ddf6-42a0-81ba-219cfd618956",
	}

	req, err := user.Unfollow(&rst)
	if err != nil {
		t.Errorf("Failed to unfollow: %v", err)
	}
	fmt.Println(req)
}

func TestGetUserFollowers(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	rst := pb.Id{
		UserId: "3a27bff7-40e4-4074-a63b-b5af91211e2f",
	}

	req, err := user.GetUserFollowers(&rst)
	if err != nil {
		t.Errorf("Failed to get followers: %v", err)
	}
	fmt.Println(req)
}

func TestGetUserFollows(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	rst := pb.Id{
		UserId: "3a27bff7-40e4-4074-a63b-b5af91211e2f",
	}

	req, err := user.GetUserFollows(&rst)
	if err != nil {
		t.Errorf("Failed to get followers: %v", err)
	}
	fmt.Println(req)
}

func TestMostPopularUser(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	rst := pb.Void{}
	req, err := user.MostPopularUser(&rst)
	if err != nil {
		t.Errorf("Failed to most popular user: %v", err)
	}
	fmt.Println(req)
}

func TestAddNationality(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	rst := pb.Nat{
		Name:        "Canadian",
		Description: "Nationality of Canada",
	}

	_, err = user.AddNationality(&rst)
	if err != nil {
		t.Errorf("Failed to add nationality: %v", err)
	}
	fmt.Println("Added nationality:", rst.Name)
}

func TestGetNationalityById(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	rst := pb.NId{
		Id: "5d57f237-bd1d-4282-8810-1887c464f7b1", // Assuming there is a nationality with ID 1
	}

	req, err := user.GetNationalityById(&rst)
	if err != nil {
		t.Errorf("Failed to get nationality: %v", err)
	}
	fmt.Println("Nationality fetched:", req)
}

func TestListNationalities(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	req := pb.Pagination{
		Limit: 10,
		Page:  1,
		Name:  "canadian", // Optional filter
	}

	res, err := user.ListNationalities(&req)
	if err != nil {
		t.Errorf("Failed to list nationalities: %v", err)
	}

	for _, nationality := range res.Nationalities {
		fmt.Println("Nationality:", nationality.Name, nationality.Description)
	}
}

func TestUpdateNationality(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	req := pb.Nationality{
		Id:          "5d57f237-bd1d-4282-8810-1887c464f7b1", // Assuming a valid ID
		Name:        "Updated Nationality Name",
		Description: "Updated Description",
	}

	_, err = user.UpdateNationality(&req)
	if err != nil {
		t.Errorf("Failed to update nationality: %v", err)
	}
	fmt.Println("Updated nationality with ID:", req.Id)
}

func TestDeleteNationality(t *testing.T) {
	db, err := ConnectUser()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	user := NewUserRepo(db)

	rst := pb.NId{
		Id: "5d57f237-bd1d-4282-8810-1887c464f7b1", // Assuming there is a nationality with ID 1
	}

	_, err = user.DeleteNationality(&rst)
	if err != nil {
		t.Errorf("Failed to delete nationality: %v", err)
	}
	fmt.Println("Deleted nationality with ID:", rst.Id)
}
