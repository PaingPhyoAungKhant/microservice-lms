// @title File Service API
// @version 1.0
// @description File management service API documentation
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
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/config"
	filePostgres "github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/infrastructure/storage"
	httpRouter "github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/interfaces/http"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/infrastructure/persistence/postgres"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
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
		"starting file service",
		zap.String("service_name", cfg.Server.ServiceName),
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)


	db, err := postgres.NewDatabase(cfg.Database)
	if err != nil {
		appLogger.Fatal("failed to create database", zap.Error(err))
	}
	appLogger.Info("database connected successfully")

	minioClient, err := storage.NewMinIOClient(cfg.MinIO, appLogger)
	if err != nil {
		appLogger.Fatal("failed to create MinIO client", zap.Error(err))
	}
	appLogger.Info("MinIO client connected successfully")

	fileRepo := filePostgres.NewPostgresFileRepository(db)

	uploadFileUseCase := usecases.NewUploadFileUseCase(fileRepo, minioClient, appLogger, cfg.Server.APIGatewayURL)
	downloadFileUseCase := usecases.NewDownloadFileUseCase(fileRepo, minioClient, appLogger)
	getFileUseCase := usecases.NewGetFileUseCase(fileRepo, appLogger, cfg.Server.APIGatewayURL)
	listFilesUseCase := usecases.NewListFilesUseCase(fileRepo, appLogger, cfg.Server.APIGatewayURL)
	deleteFileUseCase := usecases.NewDeleteFileUseCase(fileRepo, minioClient, appLogger)

	fileHttpHandler := handlers.NewFileHandler(
		uploadFileUseCase,
		downloadFileUseCase,
		getFileUseCase,
		listFilesUseCase,
		deleteFileUseCase,
		appLogger,
	)

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	httpRouter.SetUpRoutes(router, fileHttpHandler, appLogger)
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

