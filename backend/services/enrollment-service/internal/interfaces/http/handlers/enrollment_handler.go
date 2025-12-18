package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/services/enrollment-service/internal/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
)

type EnrollmentHandler struct {
	getEnrollmentUseCase         *usecases.GetEnrollmentUseCase
	findEnrollmentUseCase        *usecases.FindEnrollmentUseCase
	createEnrollmentUseCase      *usecases.CreateEnrollmentUseCase
	updateEnrollmentStatusUseCase *usecases.UpdateEnrollmentStatusUseCase
	deleteEnrollmentUseCase      *usecases.DeleteEnrollmentUseCase
	logger                       *logger.Logger
}

func NewEnrollmentHandler(
	getEnrollmentUseCase *usecases.GetEnrollmentUseCase,
	findEnrollmentUseCase *usecases.FindEnrollmentUseCase,
	createEnrollmentUseCase *usecases.CreateEnrollmentUseCase,
	updateEnrollmentStatusUseCase *usecases.UpdateEnrollmentStatusUseCase,
	deleteEnrollmentUseCase *usecases.DeleteEnrollmentUseCase,
	logger *logger.Logger,
) *EnrollmentHandler {
	return &EnrollmentHandler{
		getEnrollmentUseCase:          getEnrollmentUseCase,
		findEnrollmentUseCase:         findEnrollmentUseCase,
		createEnrollmentUseCase:       createEnrollmentUseCase,
		updateEnrollmentStatusUseCase: updateEnrollmentStatusUseCase,
		deleteEnrollmentUseCase:       deleteEnrollmentUseCase,
		logger:                        logger,
	}
}

// GetEnrollment godoc
// @Summary Get enrollment by ID
// @Description Retrieve enrollment information by enrollment ID
// @Tags enrollments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Enrollment ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Enrollment retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid enrollment ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Enrollment not found"
// @Router /{id} [get]
func (h *EnrollmentHandler) GetEnrollment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "Enrollment ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid enrollment ID format")
		return
	}

	input := usecases.GetEnrollmentInput{
		EnrollmentID: id,
	}

	output, err := h.getEnrollmentUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrEnrollmentNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

type FindEnrollmentRequest struct {
	SearchQuery      string `json:"search_query" form:"search_query" binding:"omitempty,min=3,max=255" example:"john"`
	StudentID        string `json:"student_id" form:"student_id" binding:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	CourseID         string `json:"course_id" form:"course_id" binding:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
	CourseOfferingID string `json:"course_offering_id" form:"course_offering_id" binding:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440002"`
	Status           string `json:"status" form:"status" binding:"omitempty,oneof=pending approved rejected completed" example:"pending"`
	Limit            int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=100" example:"10"`
	Offset           int    `json:"offset" form:"offset" binding:"omitempty,min=0" example:"0"`
	SortColumn       string `json:"sort_column" form:"sort_column" binding:"omitempty,oneof=created_at updated_at status" example:"created_at"`
	SortDirection    string `json:"sort_direction" form:"sort_direction" binding:"omitempty,oneof=asc desc" example:"desc"`
}

// FindEnrollment godoc
// @Summary Find enrollments with filters
// @Description Search and filter enrollments with pagination and sorting
// @Tags enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search_query query string false "Search query for student username, course name, or course offering name"
// @Param student_id query string false "Filter by student ID" Format(uuid)
// @Param course_id query string false "Filter by course ID" Format(uuid)
// @Param course_offering_id query string false "Filter by course offering ID" Format(uuid)
// @Param status query string false "Filter by status" Enums(pending, approved, rejected, completed)
// @Param limit query int false "Number of results per page" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of results to skip" default(0) minimum(0)
// @Param sort_column query string false "Column to sort by" Enums(created_at, updated_at, status)
// @Param sort_direction query string false "Sort direction" Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Enrollments retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router / [get]
func (h *EnrollmentHandler) FindEnrollment(c *gin.Context) {
	var req FindEnrollmentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	input := usecases.FindEnrollmentInput{}

	if req.SearchQuery != "" {
		input.SearchQuery = &req.SearchQuery
	}

	if req.StudentID != "" {
		input.StudentID = &req.StudentID
	}

	if req.CourseID != "" {
		input.CourseID = &req.CourseID
	}

	if req.CourseOfferingID != "" {
		input.CourseOfferingID = &req.CourseOfferingID
	}

	if req.Status != "" {
		status, err := valueobjects.NewEnrollmentStatus(req.Status)
		if err != nil {
			middleware.AbortWithError(c, http.StatusBadRequest, "Invalid status")
			return
		}
		input.Status = &status
	}

	if req.Limit > 0 {
		input.Limit = &req.Limit
	}

	if req.Offset > 0 {
		input.Offset = &req.Offset
	}

	if req.SortColumn != "" {
		input.SortColumn = &req.SortColumn
	}

	if req.SortDirection != "" {
		sortDir := repositories.SortDirection(req.SortDirection)
		input.SortDirection = &sortDir
	}

	output, err := h.findEnrollmentUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

type CreateEnrollmentRequest struct {
	StudentID          string `json:"student_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	StudentUsername    string `json:"student_username" binding:"required,min=3,max=255" example:"johndoe"`
	CourseID           string `json:"course_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
	CourseName         string `json:"course_name" binding:"required,min=1,max=255" example:"Introduction to Computer Science"`
	CourseOfferingID   string `json:"course_offering_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440002"`
	CourseOfferingName string `json:"course_offering_name" binding:"required,min=1,max=255" example:"Fall 2024"`
}

// CreateEnrollment godoc
// @Summary Create a new enrollment
// @Description Create a new enrollment for a student in a course offering
// @Tags enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateEnrollmentRequest true "Enrollment creation details"
// @Success 201 {object} map[string]interface{} "Enrollment created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "Enrollment already exists"
// @Router / [post]
func (h *EnrollmentHandler) CreateEnrollment(c *gin.Context) {
	var req CreateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	input := usecases.CreateEnrollmentInput{
		StudentID:          req.StudentID,
		StudentUsername:    req.StudentUsername,
		CourseID:           req.CourseID,
		CourseName:         req.CourseName,
		CourseOfferingID:   req.CourseOfferingID,
		CourseOfferingName: req.CourseOfferingName,
	}

	output, err := h.createEnrollmentUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrEnrollmentAlreadyExists {
			middleware.AbortWithError(c, http.StatusConflict, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, output)
}

type UpdateEnrollmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending approved rejected completed" example:"approved"`
}

// UpdateEnrollmentStatus godoc
// @Summary Update enrollment status
// @Description Update the status of an enrollment
// @Tags enrollments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Enrollment ID" Format(uuid)
// @Param request body UpdateEnrollmentStatusRequest true "Status update details"
// @Success 200 {object} map[string]interface{} "Enrollment status updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or enrollment ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Enrollment not found"
// @Router /{id}/status [put]
func (h *EnrollmentHandler) UpdateEnrollmentStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "Enrollment ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid enrollment ID format")
		return
	}

	var req UpdateEnrollmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	status, err := valueobjects.NewEnrollmentStatus(req.Status)
	if err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid status")
		return
	}

	input := usecases.UpdateEnrollmentStatusInput{
		EnrollmentID: id,
		Status:        status,
	}

	output, err := h.updateEnrollmentStatusUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrEnrollmentNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

// DeleteEnrollment godoc
// @Summary Delete enrollment
// @Description Delete an enrollment by ID
// @Tags enrollments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Enrollment ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Enrollment deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid enrollment ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Enrollment not found"
// @Router /{id} [delete]
func (h *EnrollmentHandler) DeleteEnrollment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "Enrollment ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid enrollment ID format")
		return
	}

	input := usecases.DeleteEnrollmentInput{
		EnrollmentID: id,
	}

	output, err := h.deleteEnrollmentUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrEnrollmentNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

// Health godoc
// @Summary Health check
// @Description Check if the enrollment service is running and healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service is healthy"
// @Router /health [get]
func (h *EnrollmentHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "enrollment-service",
	})
}

