// @title Course Service API
// @version 1.0
// @description Course and category management service API documentation
// @host asto-lms.local
// @BasePath /api/v1
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
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/consumer"
	appHandlers "github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/infrastructure/config"
	coursePostgres "github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/infrastructure/persistence/postgres"
	httpRouter "github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/interfaces/http"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/interfaces/http/handlers"
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
		"starting course service",
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

	categoryRepo := coursePostgres.NewPostgresCategoryRepository(db)
	courseRepo := coursePostgres.NewPostgresCourseRepository(db)
	courseCategoryRepo := coursePostgres.NewPostgresCourseCategoryRepository(db)
	offeringRepo := coursePostgres.NewPostgresCourseOfferingRepository(db)
	instructorRepo := coursePostgres.NewPostgresCourseOfferingInstructorRepository(db)
	sectionRepo := coursePostgres.NewPostgresCourseSectionRepository(db)
	moduleRepo := coursePostgres.NewPostgresSectionModuleRepository(db)

	createCategoryUseCase := usecases.NewCreateCategoryUseCase(categoryRepo, appLogger)
	findCategoryUseCase := usecases.NewFindCategoryUseCase(categoryRepo, appLogger)
	getCategoryUseCase := usecases.NewGetCategoryUseCase(categoryRepo, appLogger)
	updateCategoryUseCase := usecases.NewUpdateCategoryUseCase(categoryRepo, appLogger)
	deleteCategoryUseCase := usecases.NewDeleteCategoryUseCase(categoryRepo, appLogger)
	
	createCourseUseCase := usecases.NewCreateCourseUseCase(courseRepo, courseCategoryRepo, categoryRepo, rabbitMQ, appLogger, cfg.Server.APIGatewayURL)
	listCoursesUseCase := usecases.NewListCoursesUseCase(courseRepo, courseCategoryRepo, categoryRepo, appLogger, cfg.Server.APIGatewayURL)
	findCourseUseCase := usecases.NewFindCourseUseCase(courseRepo, courseCategoryRepo, categoryRepo, appLogger, cfg.Server.APIGatewayURL)
	getCourseUseCase := usecases.NewGetCourseUseCase(courseRepo, courseCategoryRepo, categoryRepo, appLogger, cfg.Server.APIGatewayURL)
	getCourseWithDetailsUseCase := usecases.NewGetCourseWithDetailsUseCase(
		courseRepo,
		offeringRepo,
		instructorRepo,
		sectionRepo,
		moduleRepo,
		cfg.Server.APIGatewayURL,
	)
	updateCourseUseCase := usecases.NewUpdateCourseUseCase(courseRepo, courseCategoryRepo, categoryRepo, rabbitMQ, appLogger, cfg.Server.APIGatewayURL)
	deleteCourseUseCase := usecases.NewDeleteCourseUseCase(courseRepo, rabbitMQ, appLogger)

	createOfferingUseCase := usecases.NewCreateCourseOfferingUseCase(offeringRepo, courseRepo, rabbitMQ, appLogger)
	updateOfferingUseCase := usecases.NewUpdateCourseOfferingUseCase(offeringRepo, rabbitMQ, appLogger)
	deleteOfferingUseCase := usecases.NewDeleteCourseOfferingUseCase(offeringRepo, rabbitMQ, appLogger)
	findOfferingUseCase := usecases.NewFindCourseOfferingUseCase(offeringRepo, courseRepo, appLogger)
	getOfferingUseCase := usecases.NewGetCourseOfferingUseCase(offeringRepo, courseRepo, instructorRepo, sectionRepo, moduleRepo)
	assignInstructorUseCase := usecases.NewAssignInstructorToOfferingUseCase(instructorRepo, offeringRepo, rabbitMQ, appLogger)
	removeInstructorUseCase := usecases.NewRemoveInstructorFromOfferingUseCase(instructorRepo, rabbitMQ, appLogger)

	createSectionUseCase := usecases.NewCreateCourseSectionUseCase(sectionRepo, offeringRepo, rabbitMQ, appLogger)
	updateSectionUseCase := usecases.NewUpdateCourseSectionUseCase(sectionRepo, rabbitMQ, appLogger)
	getSectionUseCase := usecases.NewGetCourseSectionUseCase(sectionRepo, appLogger)
	deleteSectionUseCase := usecases.NewDeleteCourseSectionUseCase(sectionRepo, rabbitMQ, appLogger)
	findSectionUseCase := usecases.NewFindCourseSectionUseCase(sectionRepo, appLogger)
	reorderSectionsUseCase := usecases.NewReorderCourseSectionsUseCase(sectionRepo, appLogger)

	createModuleUseCase := usecases.NewCreateSectionModuleUseCase(moduleRepo, sectionRepo, rabbitMQ, appLogger)
	updateModuleUseCase := usecases.NewUpdateSectionModuleUseCase(moduleRepo, rabbitMQ, appLogger)
	getModuleUseCase := usecases.NewGetSectionModuleUseCase(moduleRepo, appLogger)
	deleteModuleUseCase := usecases.NewDeleteSectionModuleUseCase(moduleRepo, rabbitMQ, appLogger)
	findModuleUseCase := usecases.NewFindSectionModuleUseCase(moduleRepo, appLogger)
	reorderModulesUseCase := usecases.NewReorderSectionModulesUseCase(moduleRepo, appLogger)

	categoryHandler := handlers.NewCategoryHandler(createCategoryUseCase, nil, findCategoryUseCase, getCategoryUseCase, updateCategoryUseCase, deleteCategoryUseCase, appLogger)
	courseHandler := handlers.NewCourseHandler(createCourseUseCase, listCoursesUseCase, findCourseUseCase, getCourseUseCase, getCourseWithDetailsUseCase, updateCourseUseCase, deleteCourseUseCase, appLogger)
	courseOfferingHandler := handlers.NewCourseOfferingHandler(
		createOfferingUseCase,
		updateOfferingUseCase,
		deleteOfferingUseCase,
		findOfferingUseCase,
		getOfferingUseCase,
		assignInstructorUseCase,
		removeInstructorUseCase,
		appLogger,
	)
	courseSectionHandler := handlers.NewCourseSectionHandler(
		createSectionUseCase,
		updateSectionUseCase,
		getSectionUseCase,
		deleteSectionUseCase,
		findSectionUseCase,
		reorderSectionsUseCase,
		appLogger,
	)
	sectionModuleHandler := handlers.NewSectionModuleHandler(
		createModuleUseCase,
		updateModuleUseCase,
		getModuleUseCase,
		deleteModuleUseCase,
		findModuleUseCase,
		reorderModulesUseCase,
		appLogger,
	)

	userUpdatedHandler := appHandlers.NewUserUpdatedHandler(instructorRepo, appLogger)
	zoomMeetingCreatedHandler := appHandlers.NewZoomMeetingCreatedHandler(moduleRepo, appLogger)
	eventConsumer := consumer.NewEventConsumer(rabbitMQ, userUpdatedHandler, zoomMeetingCreatedHandler, appLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := eventConsumer.Start(ctx); err != nil {
			appLogger.Error("event consumer stopped", zap.Error(err))
		}
	}()

	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	httpRouter.SetUpRoutes(
		router,
		categoryHandler,
		courseHandler,
		courseOfferingHandler,
		courseSectionHandler,
		sectionModuleHandler,
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

