package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
)

type ZoomRecordingHandler struct {
	createRecordingUseCase *usecases.CreateZoomRecordingUseCase
	updateRecordingUseCase *usecases.UpdateZoomRecordingUseCase
	getRecordingUseCase    *usecases.GetZoomRecordingUseCase
	deleteRecordingUseCase *usecases.DeleteZoomRecordingUseCase
	listRecordingsUseCase   *usecases.ListZoomRecordingsUseCase
	logger                 *logger.Logger
}

func NewZoomRecordingHandler(
	createRecordingUseCase *usecases.CreateZoomRecordingUseCase,
	updateRecordingUseCase *usecases.UpdateZoomRecordingUseCase,
	getRecordingUseCase *usecases.GetZoomRecordingUseCase,
	deleteRecordingUseCase *usecases.DeleteZoomRecordingUseCase,
	listRecordingsUseCase *usecases.ListZoomRecordingsUseCase,
	logger *logger.Logger,
) *ZoomRecordingHandler {
	return &ZoomRecordingHandler{
		createRecordingUseCase: createRecordingUseCase,
		updateRecordingUseCase: updateRecordingUseCase,
		getRecordingUseCase:    getRecordingUseCase,
		deleteRecordingUseCase: deleteRecordingUseCase,
		listRecordingsUseCase:  listRecordingsUseCase,
		logger:                 logger,
	}
}

// CreateZoomRecording godoc
// @Summary Create a new zoom recording
// @Description Create a new zoom recording record
// @Tags zoom-recordings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param recording body dtos.CreateZoomRecordingInput true "Zoom recording creation data"
// @Success 201 {object} dtos.ZoomRecordingDTO "Zoom recording created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Zoom meeting not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/recordings [post]
func (h *ZoomRecordingHandler) CreateZoomRecording(c *gin.Context) {
	var input dtos.CreateZoomRecordingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.createRecordingUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrZoomMeetingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to create zoom recording: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetZoomRecording godoc
// @Summary Get zoom recording by ID
// @Description Retrieve zoom recording information by recording ID
// @Tags zoom-recordings
// @Produce json
// @Security BearerAuth
// @Param id path string true "Zoom Recording ID" Format(uuid)
// @Success 200 {object} dtos.ZoomRecordingDTO "Zoom recording retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid recording ID"
// @Failure 404 {object} map[string]interface{} "Zoom recording not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/recordings/{id} [get]
func (h *ZoomRecordingHandler) GetZoomRecording(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "recording id is required")
		return
	}

	result, err := h.getRecordingUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrZoomRecordingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get zoom recording: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateZoomRecording godoc
// @Summary Update a zoom recording
// @Description Update an existing zoom recording
// @Tags zoom-recordings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Zoom Recording ID" Format(uuid)
// @Param recording body dtos.UpdateZoomRecordingInput true "Zoom recording update data"
// @Success 200 {object} dtos.ZoomRecordingDTO "Zoom recording updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Zoom recording not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/recordings/{id} [put]
func (h *ZoomRecordingHandler) UpdateZoomRecording(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "recording id is required")
		return
	}

	var input dtos.UpdateZoomRecordingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.updateRecordingUseCase.Execute(c.Request.Context(), id, input)
	if err != nil {
		if err == usecases.ErrZoomRecordingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to update zoom recording: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteZoomRecording godoc
// @Summary Delete a zoom recording
// @Description Delete a zoom recording by its ID
// @Tags zoom-recordings
// @Produce json
// @Security BearerAuth
// @Param id path string true "Zoom Recording ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Zoom recording deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid recording ID"
// @Failure 404 {object} map[string]interface{} "Zoom recording not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/recordings/{id} [delete]
func (h *ZoomRecordingHandler) DeleteZoomRecording(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "recording id is required")
		return
	}

	err := h.deleteRecordingUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrZoomRecordingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to delete zoom recording: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "zoom recording deleted successfully"})
}

// ListZoomRecordings godoc
// @Summary List zoom recordings for a meeting
// @Description Retrieve all zoom recordings for a specific zoom meeting
// @Tags zoom-recordings
// @Produce json
// @Security BearerAuth
// @Param meeting_id path string true "Zoom Meeting ID" Format(uuid)
// @Success 200 {array} dtos.ZoomRecordingDTO "Zoom recordings retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid meeting ID"
// @Failure 404 {object} map[string]interface{} "Zoom meeting not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/recordings/meeting/{meeting_id} [get]
func (h *ZoomRecordingHandler) ListZoomRecordings(c *gin.Context) {
	meetingID := c.Param("meeting_id")
	if meetingID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "meeting id is required")
		return
	}

	result, err := h.listRecordingsUseCase.Execute(c.Request.Context(), meetingID)
	if err != nil {
		if err == usecases.ErrZoomMeetingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to list zoom recordings: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

