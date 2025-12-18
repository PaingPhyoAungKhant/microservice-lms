// @title Zoom Service API
// @version 1.0
// @description Zoom meeting and recording management service API documentation
// @host asto-lms.local
// @BasePath /api/v1
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
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/config"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/external/zoom"
	zoomPostgres "github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/infrastructure/persistence/postgres"
	httpRouter "github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/interfaces/http"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/interfaces/http/handlers"
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

	appLogger, err := logger.NewLogger(cfg.Server.Environment)
	if err != nil {
		log.Fatal("failed to create logger: %w", err)
	}
	defer appLogger.Core().Sync()

	appLogger.Info(fmt.Sprintf("config: %+v", cfg))
	appLogger.Info(
		"starting zoom service",
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

	meetingRepo := zoomPostgres.NewPostgresZoomMeetingRepository(db)
	recordingRepo := zoomPostgres.NewPostgresZoomRecordingRepository(db)

	zoomClient := zoom.NewZoomClient(&cfg.Zoom)

	userID := os.Getenv("ZOOM_USER_ID")
	if userID == "" {
		appLogger.Fatal("ZOOM_USER_ID environment variable is required")
	}

	createMeetingUseCase := usecases.NewCreateZoomMeetingUseCase(meetingRepo, zoomClient, rabbitMQ, appLogger, userID)
	updateMeetingUseCase := usecases.NewUpdateZoomMeetingUseCase(meetingRepo, zoomClient)
	getMeetingUseCase := usecases.NewGetZoomMeetingUseCase(meetingRepo)
	deleteMeetingUseCase := usecases.NewDeleteZoomMeetingUseCase(meetingRepo, zoomClient)
	getMeetingByModuleUseCase := usecases.NewGetZoomMeetingByModuleUseCase(meetingRepo)

	createRecordingUseCase := usecases.NewCreateZoomRecordingUseCase(recordingRepo, meetingRepo)
	updateRecordingUseCase := usecases.NewUpdateZoomRecordingUseCase(recordingRepo)
	getRecordingUseCase := usecases.NewGetZoomRecordingUseCase(recordingRepo)
	deleteRecordingUseCase := usecases.NewDeleteZoomRecordingUseCase(recordingRepo)
	listRecordingsUseCase := usecases.NewListZoomRecordingsUseCase(recordingRepo, meetingRepo)

	zoomMeetingHandler := handlers.NewZoomMeetingHandler(
		createMeetingUseCase,
		updateMeetingUseCase,
		getMeetingUseCase,
		deleteMeetingUseCase,
		getMeetingByModuleUseCase,
		appLogger,
	)

	zoomRecordingHandler := handlers.NewZoomRecordingHandler(
		createRecordingUseCase,
		updateRecordingUseCase,
		getRecordingUseCase,
		deleteRecordingUseCase,
		listRecordingsUseCase,
		appLogger,
	)

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	httpRouter.SetUpRoutes(
		router,
		zoomMeetingHandler,
		zoomRecordingHandler,
		appLogger,
	)

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
	if err := utils.GracefulShutDown(server, nil, db, appLogger); err != nil {
		appLogger.Error("failed to shutdown server", zap.Error(err))
	}
	appLogger.Info("Server shutdown completed")
}

