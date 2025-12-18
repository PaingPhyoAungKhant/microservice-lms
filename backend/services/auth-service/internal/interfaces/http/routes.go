package http

import (
	"github.com/gin-gonic/gin"
	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
	_ "github.com/paingphyoaungkhant/asto-microservice/services/auth-service/cmd/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, logger *logger.Logger) {

	// Middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())
	router.GET("/health", authHandler.Health)

	// Swagger documentation
	router.GET("/api/v1/auth/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Routes
	api := router.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterStudent)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.GET("/verify", authHandler.Verify)
		auth.GET("/verify-email", authHandler.VerifyEmail)
		auth.POST("/request-email-verify", authHandler.RequestEmailVerify)
		auth.POST("/refresh-token", authHandler.RefreshToken)
	}
}