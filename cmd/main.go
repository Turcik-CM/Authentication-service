package main

import (
	pb "auth-service/genproto/user"
	"auth-service/pkg/config"
	"auth-service/pkg/logs"
	"auth-service/service"
	"auth-service/storage/postgres"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger := logs.InitLogger()
	cfg := config.Load()

	db, err := postgres.ConnectPostgres(cfg)
	if err != nil {
		logger.Error("Error connecting to database")
		log.Fatal(err)
	}
	defer db.Close()

	authst := postgres.NewAuthRepo(db)
	userst := postgres.NewUserRepo(db)

	authSr := service.NewAuthService(authst, userst, logger)

	listen, err := net.Listen("tcp", cfg.USER_PORT)
	fmt.Println("listening on port " + cfg.USER_PORT)
	if err != nil {
		logger.Error("Error listening on port " + cfg.USER_PORT)
		log.Fatal(err)
	}

	service := grpc.NewServer()
	pb.RegisterUserServiceServer(service, authSr)

	if err := service.Serve(listen); err != nil {
		logger.Error("Error starting server")
		log.Fatal(err)
	}
}
