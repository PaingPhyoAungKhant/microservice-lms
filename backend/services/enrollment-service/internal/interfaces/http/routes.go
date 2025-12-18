package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/cmd/docs"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(router *gin.Engine, handler *handlers.EnrollmentHandler, logger *logger.Logger) {
	// Middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Routes
	router.GET("/health", handler.Health)

	// Swagger documentation
	router.GET("/api/v1/enrollments/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Enrollment Routes
	api := router.Group("/api/v1")
	{
		enrollmentRouter := api.Group("/enrollments")
		enrollmentRouter.POST("", handler.CreateEnrollment)
		enrollmentRouter.GET("/:id", handler.GetEnrollment)
		enrollmentRouter.PUT("/:id/status", handler.UpdateEnrollmentStatus)
		enrollmentRouter.DELETE("/:id", handler.DeleteEnrollment)
		enrollmentRouter.GET("", handler.FindEnrollment)
	}
}

