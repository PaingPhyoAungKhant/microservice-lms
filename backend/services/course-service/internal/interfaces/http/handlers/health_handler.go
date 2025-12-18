package handlers

import (
	"github.com/gin-gonic/gin"
)

// Health godoc
// @Summary Health check
// @Description Check the health status of the course service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "Service is healthy"
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "course-service",
	})
}

