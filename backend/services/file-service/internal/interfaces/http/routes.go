package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/paingphyoaungkhant/asto-microservice/services/file-service/cmd/docs"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(router *gin.Engine, handler *handlers.FileHandler, logger *logger.Logger) {
	// Middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	router.GET("/health", handler.Health)

	// Swagger documentation
	router.GET("/api/v1/files/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// File Routes
	api := router.Group("/api/v1")
	{
		fileRouter := api.Group("/files")
		fileRouter.POST("", handler.UploadFile)
		fileRouter.GET("/:id", handler.GetFile)
		fileRouter.GET("/:id/download", handler.DownloadFile)
		fileRouter.GET("", handler.ListFiles)
		fileRouter.DELETE("/:id", handler.DeleteFile)
		
		bucketRouter := api.Group("/buckets")
		{
			bucketRouter.GET("/:bucket/files/:id/download", handler.DownloadFileByBucket)
		}
	}
}

