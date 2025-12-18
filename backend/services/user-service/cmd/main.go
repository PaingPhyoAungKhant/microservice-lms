// @title User Service API
// @version 1.0
// @description User management service API documentation
// @host asto-lms.local
// @BasePath /api/v1/users
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/infrastructure/config"
	httpRouter "github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/interfaces/http"
	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: %w", err)
	}

	fmt.Println("config: %+v", config)

	appLogger, err := logger.NewLogger(config.Server.Environment)
	if err != nil {
		log.Fatal("failed to create logger: %w", err)
	}
	defer appLogger.Core().Sync()
	appLogger.Info(
		"starting user service",
		zap.String("service_name", config.Server.ServiceName),
		zap.String("environment", config.Server.Environment),
		zap.String("port", config.Server.Port),
	)

	db, err := postgres.NewDatabase(config.Database)
	if err != nil {
		appLogger.Fatal("failed to create database", zap.Error(err))
	}
	appLogger.Info("database connected successfully")


	rabbitMQ, err := messaging.NewRabbitMQ(&config.RabbitMQ, appLogger)
	if err != nil {
		appLogger.Error("failed to create rabbitmq", zap.Error(err))
	} else {
		appLogger.Info("rabbitmq connected successfully")
	}

	redis, err := utils.NewRedis(&config.Redis)
	if err != nil {
		appLogger.Error("failed to create redis", zap.Error(err))
	} else {
		appLogger.Info("redis connected successfully")
	}

	userRepo := postgres.NewPostgresUserRepository(db)
	createUserUseCase := usecases.NewCreateUserUseCase(userRepo, rabbitMQ, appLogger, redis, config.Server.APIGatewayURL)
	updateUserUseCase := usecases.NewUpdateUserUseCase(userRepo, rabbitMQ, appLogger)
	getUserUseCase := usecases.NewGetUserUseCase(userRepo, appLogger)
	findUserUseCase := usecases.NewFindUserUseCase(userRepo, appLogger)
	deleteUserUseCase := usecases.NewDeleteUserUseCase(userRepo, rabbitMQ, appLogger)

	userHttpHandler := handlers.NewUserHandler(createUserUseCase, getUserUseCase, updateUserUseCase, findUserUseCase, deleteUserUseCase, appLogger)
	if config.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	httpRouter.SetUpRoutes(router, userHttpHandler, appLogger)
	server := &http.Server{
		Addr:         ":" + config.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  16 * time.Second,
	}

	go func() {
		appLogger.Info("starting server", zap.String("port", config.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")
	if err := utils.GracefulShutDown(server, rabbitMQ, db, appLogger); err != nil {
		appLogger.Error("failed to shutdown server", zap.Error(err))
	}
	appLogger.Info("Server shutdown completed")

}
