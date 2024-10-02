package postgres

import (
	pb "auth-service/genproto/user"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"testing"
)

func Connect() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "postgres", "dodi", "auth_tw")

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestRegister(t *testing.T) {

	db, err := Connect()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}

	rst := pb.RegisterRequest{
		Email:     "dodi",
		Phone:     "dodi",
		FirstName: "dodi",
		LastName:  "dodi",
		Username:  "dodi",
		Country:   "Uzbekistan",
		Password:  "dodi",
		Bio:       "-----------------------------",
	}

	auth := NewAuthRepo(db)

	req, err := auth.Register(&rst)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}

	fmt.Println(req)

}

func TestLoginEmail(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}

	rst := pb.LoginEmailRequest{
		Email:    "dodi",
		Password: "dodi",
	}

	auth := NewAuthRepo(db)

	req, err := auth.LoginEmail(&rst)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}

	fmt.Println(req)
}

func TestLoginUsername(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	rst := pb.LoginUsernameRequest{
		Username: "dodi",
		Password: "dodi",
	}
	auth := NewAuthRepo(db)
	req, err := auth.LoginUsername(&rst)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}
	fmt.Println(req)
}

func TestGetUserByEmail(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	res := pb.Email{
		Email: "dodi",
	}
	auth := NewAuthRepo(db)
	req, err := auth.GetUserByEmail(&res)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}
	fmt.Println(req)
}

func TestRegisterAdmin(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	rst := pb.Message{
		Message: "123321",
	}
	auth := NewAuthRepo(db)
	req, err := auth.RegisterAdmin(&rst)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}
	fmt.Println(req)
}

func TestUpdatePassword(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	rst := pb.UpdatePasswordReq{
		Id:       "cdb90f0d-c69d-432a-b8a0-6fc40e283ccb",
		Password: "12221",
	}
	auth := NewAuthRepo(db)
	req, err := auth.UpdatePassword(&rst)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}
	fmt.Println(req)
}
