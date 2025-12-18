// @title Enrollment Service API
// @version 1.0
// @description Enrollment management service API documentation
// @host asto-lms.local
// @BasePath /api/v1/enrollments
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/consumer"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/infrastructure/config"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/infrastructure/persistence/postgres"
	httpRouter "github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/interfaces/http"
	enrollmentHandlers "github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/interfaces/http/handlers"
	sharedPostgres "github.com/paingphyoaungkhant/asto-microservice/shared/infrastructure/persistence/postgres"
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
		"starting enrollment service",
		zap.String("service_name", cfg.Server.ServiceName),
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	db, err := sharedPostgres.NewDatabase(cfg.Database)
	if err != nil {
		appLogger.Fatal("failed to create database", zap.Error(err))
	}
	appLogger.Info("database connected successfully")

	rabbitMQ, err := messaging.NewRabbitMQ(&cfg.RabbitMQ, appLogger)
	if err != nil {
		appLogger.Error("failed to create rabbitmq", zap.Error(err))
	} else {
		appLogger.Info("rabbitmq connected successfully")
		defer rabbitMQ.Close()
	}

	enrollmentRepo := postgres.NewPostgresEnrollmentRepository(db)


	userUpdatedHandler := handlers.NewUserUpdatedHandler(enrollmentRepo, appLogger)
	courseUpdatedHandler := handlers.NewCourseUpdatedHandler(enrollmentRepo, appLogger)
	courseOfferingUpdatedHandler := handlers.NewCourseOfferingUpdatedHandler(enrollmentRepo, appLogger)


	eventConsumer := consumer.NewEventConsumer(
		rabbitMQ,
		userUpdatedHandler,
		courseUpdatedHandler,
		courseOfferingUpdatedHandler,
		appLogger,
	)

	createEnrollmentUseCase := usecases.NewCreateEnrollmentUseCase(enrollmentRepo, rabbitMQ, appLogger)
	getEnrollmentUseCase := usecases.NewGetEnrollmentUseCase(enrollmentRepo, appLogger)
	findEnrollmentUseCase := usecases.NewFindEnrollmentUseCase(enrollmentRepo, appLogger)
	updateEnrollmentStatusUseCase := usecases.NewUpdateEnrollmentStatusUseCase(enrollmentRepo, rabbitMQ, appLogger)
	deleteEnrollmentUseCase := usecases.NewDeleteEnrollmentUseCase(enrollmentRepo, rabbitMQ, appLogger)


	enrollmentHttpHandler := enrollmentHandlers.NewEnrollmentHandler(
		getEnrollmentUseCase,
		findEnrollmentUseCase,
		createEnrollmentUseCase,
		updateEnrollmentStatusUseCase,
		deleteEnrollmentUseCase,
		appLogger,
	)

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	httpRouter.SetUpRoutes(router, enrollmentHttpHandler, appLogger)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  16 * time.Second,
	}


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := eventConsumer.Start(ctx); err != nil {
			appLogger.Error("event consumer stopped", zap.Error(err))
		}
	}()


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
	cancel() 

	if err := utils.GracefulShutDown(server, rabbitMQ, db, appLogger); err != nil {
		appLogger.Error("failed to shutdown server", zap.Error(err))
	}
	appLogger.Info("Server shutdown completed")
}

