package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"week_06_web/auth"
	"week_06_web/handler"
	"week_06_web/middleware"
	"week_06_web/service"
)

func New(s *service.UserService, tokens *auth.TokenManager) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.RequestID(),
		middleware.AccessLog(slog.Default()),
		middleware.CORS(),
	)

	userHandler := handler.NewUserHandler(s, tokens)

	api := r.Group("/api/v1")
	{
		api.POST("/users/register", userHandler.Register)
		api.POST("/users/login", userHandler.Login)
		api.POST("/users/logout", userHandler.Logout)
		api.GET("/users/:id", userHandler.GetByID)

		protected := api.Group("")
		protected.Use(middleware.JWTAuth(tokens))
		{
			protected.GET("/users/me", userHandler.Me)
			protected.PUT("/users/:id/password", userHandler.ChangePassword)
		}
	}
	return r
}
