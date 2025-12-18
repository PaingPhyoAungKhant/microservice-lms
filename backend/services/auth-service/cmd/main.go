// @title Auth Service API
// @version 1.0
// @description Authentication and authorization service API documentation
// @host asto-lms.local
// @BasePath /api/v1/auth
// @schemes http
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
	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/infrastructure/config"
	httpRouter "github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/interfaces/http"
	"github.com/paingphyoaungkhant/asto-microservice/services/auth-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"github.com/paingphyoaungkhant/asto-microservice/shared/utils"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: %w", err)
	}

	fmt.Printf("config: %+v\n", cfg)

	appLogger, err := logger.NewLogger(cfg.Server.Environment)
	if err != nil {
		log.Fatal("failed to create logger: %w", err)
	}
	defer appLogger.Core().Sync()
	appLogger.Info(
		"starting auth service",
		zap.String("service_name", cfg.Server.ServiceName),
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	db, err := postgres.NewDatabase(cfg.Database)
	if err != nil {
		appLogger.Fatal("failed to create database", zap.Error(err))
	}
	appLogger.Info("database connected successfully")

	rabbitMQ, err := messaging.NewRabbitMQ(&cfg.RabbitMQ, appLogger)
	if err != nil {
		appLogger.Error("failed to create rabbitmq", zap.Error(err))
	} else {
		appLogger.Info("rabbitmq connected successfully")
	}

	redis, err := utils.NewRedis(&cfg.Redis)
	if err != nil {
		appLogger.Error("failed to create redis", zap.Error(err))
	} else {
		appLogger.Info("redis connected successfully")
	}

	jwtManager := utils.NewJwtManager(
		cfg.Jwt.SecretKey,
		cfg.Jwt.AccessTokenDuration,
		cfg.Jwt.RefreshTokenDuration,
	)

	userRepo := postgres.NewPostgresUserRepository(db)

	loginUseCase := usecases.NewLoginUseCase(userRepo, rabbitMQ, appLogger, jwtManager, redis)
	registerStudentUseCase := usecases.NewRegisterStudentUseCase(userRepo, rabbitMQ, appLogger, &cfg.RabbitMQ, redis, cfg.Server.APIGatewayURL)
	forgotPasswordUseCase := usecases.NewForgotPasswordUseCase(userRepo, rabbitMQ, appLogger)
	verifyOTPUseCase := usecases.NewVerifyOTPUseCase(userRepo, appLogger, redis)
	resetPasswordUseCase := usecases.NewResetPasswordUseCase(userRepo, rabbitMQ, appLogger, redis)
	verifyUseCase := usecases.NewVerifyUseCase(userRepo, rabbitMQ, appLogger, jwtManager, redis)
	verifyEmailUseCase := usecases.NewVerifyEmailUseCase(userRepo, rabbitMQ, appLogger, redis)
	requestEmailVerifyUseCase := usecases.NewRequestEmailVerifyUseCase(userRepo, rabbitMQ, appLogger, redis, cfg.Server.APIGatewayURL)
	refreshTokenUseCase := usecases.NewRefreshTokenUseCase(userRepo, appLogger, jwtManager, redis)

	authHandler := handlers.NewAuthHandler(
		loginUseCase,
		registerStudentUseCase,
		forgotPasswordUseCase,
		resetPasswordUseCase,
		verifyOTPUseCase,
		verifyUseCase,
		verifyEmailUseCase,
		requestEmailVerifyUseCase,
		refreshTokenUseCase,
	)

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	httpRouter.SetupRoutes(router, authHandler, appLogger)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  16 * time.Second,
	}

	go func() {
		appLogger.Info("starting server", zap.String("port", cfg.Server.Port))
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
