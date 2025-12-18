package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Message string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			requestID := GetRequestID(c)
			
		
			log.Error(
				"Request error",
				zap.Error(err.Err),
			)
			
			
			status := c.Writer.Status()
			if status == http.StatusOK {
				status = http.StatusInternalServerError
			}
			
			c.JSON(status, ErrorResponse{
				Error:     http.StatusText(status),
				Message:   err.Error(),
				RequestID: requestID,
			})
		}
	}
}

func AbortWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Error: "error",
		Message: message,
		RequestID: GetRequestID(c),
	})
	c.Abort()
}

func AbortWithErrorDetails(c *gin.Context, code int, message string, details map[string]interface{}) {
	c.JSON(code, ErrorResponse{
		Error:     http.StatusText(code),
		Message:   message,
		RequestID: GetRequestID(c),
		Details:   details,
	})
	c.Abort()
}


const RequestIDKey = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
	
		requestID := c.GetHeader(RequestIDKey)
		
		
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		
		c.Set("request_id", requestID)
		c.Header(RequestIDKey, requestID)
		
		c.Next()
	}
}

 
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}