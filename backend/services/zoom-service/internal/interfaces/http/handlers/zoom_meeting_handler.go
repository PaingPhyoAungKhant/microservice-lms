package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/zoom-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
)

type ZoomMeetingHandler struct {
	createMeetingUseCase      *usecases.CreateZoomMeetingUseCase
	updateMeetingUseCase     *usecases.UpdateZoomMeetingUseCase
	getMeetingUseCase         *usecases.GetZoomMeetingUseCase
	deleteMeetingUseCase     *usecases.DeleteZoomMeetingUseCase
	getMeetingByModuleUseCase *usecases.GetZoomMeetingByModuleUseCase
	logger                    *logger.Logger
}

func NewZoomMeetingHandler(
	createMeetingUseCase *usecases.CreateZoomMeetingUseCase,
	updateMeetingUseCase *usecases.UpdateZoomMeetingUseCase,
	getMeetingUseCase *usecases.GetZoomMeetingUseCase,
	deleteMeetingUseCase *usecases.DeleteZoomMeetingUseCase,
	getMeetingByModuleUseCase *usecases.GetZoomMeetingByModuleUseCase,
	logger *logger.Logger,
) *ZoomMeetingHandler {
	return &ZoomMeetingHandler{
		createMeetingUseCase:      createMeetingUseCase,
		updateMeetingUseCase:     updateMeetingUseCase,
		getMeetingUseCase:         getMeetingUseCase,
		deleteMeetingUseCase:     deleteMeetingUseCase,
		getMeetingByModuleUseCase: getMeetingByModuleUseCase,
		logger:                    logger,
	}
}

// CreateZoomMeeting godoc
// @Summary Create a new zoom meeting
// @Description Create a new zoom meeting for a section module
// @Tags zoom-meetings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param meeting body dtos.CreateZoomMeetingInput true "Zoom meeting creation data"
// @Success 201 {object} dtos.ZoomMeetingDTO "Zoom meeting created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 409 {object} map[string]interface{} "Zoom meeting already exists for this section module"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/meetings [post]
func (h *ZoomMeetingHandler) CreateZoomMeeting(c *gin.Context) {
	var input dtos.CreateZoomMeetingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.createMeetingUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrZoomMeetingAlreadyExists {
			middleware.AbortWithError(c, http.StatusConflict, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to create zoom meeting: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetZoomMeeting godoc
// @Summary Get zoom meeting by ID
// @Description Retrieve zoom meeting information by meeting ID
// @Tags zoom-meetings
// @Produce json
// @Security BearerAuth
// @Param id path string true "Zoom Meeting ID" Format(uuid)
// @Success 200 {object} dtos.ZoomMeetingDTO "Zoom meeting retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid meeting ID"
// @Failure 404 {object} map[string]interface{} "Zoom meeting not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/meetings/{id} [get]
func (h *ZoomMeetingHandler) GetZoomMeeting(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "meeting id is required")
		return
	}

	result, err := h.getMeetingUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrZoomMeetingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get zoom meeting: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateZoomMeeting godoc
// @Summary Update a zoom meeting
// @Description Update an existing zoom meeting
// @Tags zoom-meetings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Zoom Meeting ID" Format(uuid)
// @Param meeting body dtos.UpdateZoomMeetingInput true "Zoom meeting update data"
// @Success 200 {object} dtos.ZoomMeetingDTO "Zoom meeting updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Zoom meeting not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/meetings/{id} [put]
func (h *ZoomMeetingHandler) UpdateZoomMeeting(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "meeting id is required")
		return
	}

	var input dtos.UpdateZoomMeetingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.updateMeetingUseCase.Execute(c.Request.Context(), id, input)
	if err != nil {
		if err == usecases.ErrZoomMeetingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to update zoom meeting: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteZoomMeeting godoc
// @Summary Delete a zoom meeting
// @Description Delete a zoom meeting by its ID
// @Tags zoom-meetings
// @Produce json
// @Security BearerAuth
// @Param id path string true "Zoom Meeting ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Zoom meeting deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid meeting ID"
// @Failure 404 {object} map[string]interface{} "Zoom meeting not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/meetings/{id} [delete]
func (h *ZoomMeetingHandler) DeleteZoomMeeting(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "meeting id is required")
		return
	}

	err := h.deleteMeetingUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrZoomMeetingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to delete zoom meeting: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "zoom meeting deleted successfully"})
}

// GetZoomMeetingByModule godoc
// @Summary Get zoom meeting by section module ID
// @Description Retrieve zoom meeting information by section module ID
// @Tags zoom-meetings
// @Produce json
// @Security BearerAuth
// @Param module_id path string true "Section Module ID" Format(uuid)
// @Success 200 {object} dtos.ZoomMeetingDTO "Zoom meeting retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid module ID"
// @Failure 404 {object} map[string]interface{} "Zoom meeting not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /zoom/meetings/module/{module_id} [get]
func (h *ZoomMeetingHandler) GetZoomMeetingByModule(c *gin.Context) {
	moduleID := c.Param("module_id")
	if moduleID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "module id is required")
		return
	}

	result, err := h.getMeetingByModuleUseCase.Execute(c.Request.Context(), moduleID)
	if err != nil {
		if err == usecases.ErrZoomMeetingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get zoom meeting: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

