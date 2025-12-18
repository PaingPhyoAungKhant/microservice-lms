package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)


func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		
		c.Next()
		
		latency := time.Since(start)
		
		requestID := GetRequestID(c)
		
		fields := map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}
		
		if userID, exists := c.Get("user_id"); exists {
			fields["user_id"] = userID
		}
		
		if c.Writer.Status() >= 500 {
			log.Error("Server error", zap.Any("fields", fields))
		} else if c.Writer.Status() >= 400 {
			log.Warn("Client error", zap.Any("fields", fields))
		} else {
			log.Info("Request completed", zap.Any("fields", fields))
		}
	}
}

