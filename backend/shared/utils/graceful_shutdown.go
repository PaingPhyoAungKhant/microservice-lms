package utils

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/messaging"
	"go.uber.org/zap"
)

func GracefulShutDown(server *http.Server, rabbitMQ *messaging.RabbitMQ, db *sql.DB, logger *logger.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)

	defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        return fmt.Errorf("server forced to shutdown: %w", err)
    }
    logger.Info("Server Shutdown Successfully")
	if rabbitMQ != nil {
		if err := rabbitMQ.Close(); err != nil {
			logger.Error("Error while closing RabbitMQ: %w", zap.Error(err))
		}else {
			logger.Info("RabbitMQ closed successfully")
		}
	}

    if db != nil {
        if err := db.Close(); err != nil {
            return fmt.Errorf("error while closing database: %w", err)
        }
		logger.Info("Database closed successfully")
	}
	return nil
}