package main

import (
	"auth-service/api"
	"auth-service/api/handler"
	"auth-service/genproto/user"
	"auth-service/pkg/config"
	"auth-service/pkg/logs"
	"auth-service/service"
	"auth-service/storage/postgres"
	"auth-service/storage/redis"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger := logs.InitLogger()
	cofg := config.Load()
	redisClient := redis.ConnectDB()

	db, err := postgres.ConnectPostgres(cofg)
	if err != nil {
		logger.Error("Error connecting to database", "error", err)
		log.Fatal(err)
	}

	userSt := postgres.NewUserRepo(db)
	userSr := service.NewUserService(userSt, logger)
	listen, err := net.Listen("tcp", cofg.USER_PORT)
	if err != nil {
		logger.Error("Error listening on port "+cofg.USER_PORT, "error", err)
		log.Fatal(err)
	}

	go func() {
		server := grpc.NewServer()
		user.RegisterUserServiceServer(server, userSr)
		logger.Info("Starting server on port " + cofg.USER_PORT)
		log.Println("Starting server on port " + cofg.USER_PORT)

		if err := server.Serve(listen); err != nil {
			logger.Error("Error starting server on port "+cofg.USER_PORT, "error", err)
			log.Fatal(err)
		}
	}()

	//------------------------------------------------------------------

	authSt := postgres.NewAuthRepo(db)
	authSr := service.NewAuthService(authSt, logger)
	redis1 := redis.NewRedisStorage(redisClient, logger)
	authHd := handler.NewAuthHandler(logger, authSr, redis1)
	router := api.NewRouter(authHd)

	router.InitRouter()
	err = router.RunRouter(cofg)

	if err != nil {
		logger.Error("Error starting server", "error", err)
		log.Fatal(err)
	}
}
