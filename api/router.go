package api

import (
	"auth-service/api/handler"
	"auth-service/pkg/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth-service/api/docs"
)

type Router interface {
	InitRouter()
	RunRouter(cfg config.Config) error
}

type router struct {
	router   *gin.Engine
	handlers handler.AuthHandler
}

func NewRouter(authHandler handler.AuthHandler) Router {
	routers := gin.Default()
	return &router{
		router:   routers,
		handlers: authHandler,
	}
}

// @title Authenfication service
// @version 1.0
// @description server for siginIn or signUp
// @BasePath /auth
// @schemes http
func (r *router) InitRouter() {

	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.router.Group("/auth")
	{
		auth.POST("/register", r.handlers.Register)
		auth.POST("/login/email", r.handlers.LoginEmail)
		auth.POST("/login/username", r.handlers.LoginUsername)
		auth.POST("/accept-code", r.handlers.AcceptCodeToRegister)
		auth.POST("/forgot-password", r.handlers.ForgotPassword)
		auth.POST("/register-admin", r.handlers.RegisterAdmin)
		auth.POST("/reset-password", r.handlers.ResetPassword)
	}
}

func (r *router) RunRouter(cfg config.Config) error {
	return r.router.Run(cfg.AUTH_PORT)
}
