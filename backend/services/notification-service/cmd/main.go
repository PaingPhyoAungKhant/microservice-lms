package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/application/consumer"
	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/application/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/domain/templates"
	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/infrastructure/config"
	"github.com/paingphyoaungkhant/asto-microservice/services/notification-service/internal/infrastructure/email"
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
		"starting notification service",
		zap.String("service_name", cfg.Server.ServiceName),
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	
	if len(templates.EmailVerification) == 0 {
		appLogger.Fatal("email verification template is empty - check embed directive")
	}
	appLogger.Info("email templates loaded",
		zap.Int("email_verification_template_length", len(templates.EmailVerification)),
		zap.Int("forgot_password_template_length", len(templates.ForgotPasswordOTP)),
	)

	rabbitMQ, err := messaging.NewRabbitMQ(&cfg.RabbitMQ, appLogger)
	if err != nil {
		appLogger.Fatal("failed to create rabbitmq", zap.Error(err))
	}
	defer rabbitMQ.Close()
	appLogger.Info("rabbitmq connected successfully")

	redis, err := utils.NewRedis(&cfg.Redis)
	if err != nil {
		appLogger.Fatal("failed to create redis", zap.Error(err))
	}
	appLogger.Info("redis connected successfully")

	emailService := email.NewEmailService(&cfg.SMTP, appLogger)

	userCreatedHandler := handlers.NewUserCreatedHandler(emailService, appLogger)
	studentRegisteredHandler := handlers.NewStudentRegisteredHandler(emailService, appLogger)
	emailVerificationRequestHandler := handlers.NewEmailVerificationRequestHandler(emailService, appLogger)
	forgotPasswordHandler := handlers.NewForgotPasswordHandler(emailService, redis, appLogger)

	eventConsumer := consumer.NewEventConsumer(
		rabbitMQ,
		userCreatedHandler,
		studentRegisteredHandler,
		emailVerificationRequestHandler,
		forgotPasswordHandler,
		appLogger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := eventConsumer.Start(ctx); err != nil {
			appLogger.Error("event consumer stopped", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down notification service...")
	cancel()
	appLogger.Info("Notification service shutdown completed")
}

