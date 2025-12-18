package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/cmd/docs"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/interfaces/http/handlers"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(
	router *gin.Engine,
	zoomMeetingHandler *handlers.ZoomMeetingHandler,
	zoomRecordingHandler *handlers.ZoomRecordingHandler,
	logger *logger.Logger,
) {
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.Health)

	router.GET("/api/v1/zoom/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := router.Group("/api/v1")
	{
		zoomRoutes := api.Group("/zoom")
		{
			meetingRoutes := zoomRoutes.Group("/meetings")
			{
				meetingRoutes.POST("", zoomMeetingHandler.CreateZoomMeeting)
				meetingRoutes.GET("/:id", zoomMeetingHandler.GetZoomMeeting)
				meetingRoutes.PUT("/:id", zoomMeetingHandler.UpdateZoomMeeting)
				meetingRoutes.DELETE("/:id", zoomMeetingHandler.DeleteZoomMeeting)
				meetingRoutes.GET("/module/:module_id", zoomMeetingHandler.GetZoomMeetingByModule)
			}

			recordingRoutes := zoomRoutes.Group("/recordings")
			{
				recordingRoutes.POST("", zoomRecordingHandler.CreateZoomRecording)
				recordingRoutes.GET("/:id", zoomRecordingHandler.GetZoomRecording)
				recordingRoutes.PUT("/:id", zoomRecordingHandler.UpdateZoomRecording)
				recordingRoutes.DELETE("/:id", zoomRecordingHandler.DeleteZoomRecording)
				recordingRoutes.GET("/meeting/:meeting_id", zoomRecordingHandler.ListZoomRecordings)
			}
		}
	}
}

