package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/paingphyoaungkhant/asto-microservice/services/user-service/cmd/docs"
	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func SetUpRoutes(router *gin.Engine, handler *handlers.UserHandler, logger *logger.Logger){
	// Middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Routes
	router.GET("/health", handler.Health)

	// Swagger documentation
	router.GET("/api/v1/users/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// User Routes
	api := router.Group("/api/v1")
	{
		userRouter := api.Group("/users")
		userRouter.POST("", handler.CreateUser)
		userRouter.GET("/:id", handler.GetUser)
		userRouter.PUT("/:id", handler.UpdateUser)
		userRouter.DELETE("/:id", handler.DeleteUser)
		userRouter.GET("", handler.FindUser)
	}

}