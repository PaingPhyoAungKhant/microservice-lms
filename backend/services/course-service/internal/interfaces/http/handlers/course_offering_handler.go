package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
)

type CourseOfferingHandler struct {
	createOfferingUseCase      *usecases.CreateCourseOfferingUseCase
	updateOfferingUseCase      *usecases.UpdateCourseOfferingUseCase
	deleteOfferingUseCase      *usecases.DeleteCourseOfferingUseCase
	findOfferingUseCase        *usecases.FindCourseOfferingUseCase
	getOfferingUseCase         *usecases.GetCourseOfferingUseCase
	assignInstructorUseCase    *usecases.AssignInstructorToOfferingUseCase
	removeInstructorUseCase   *usecases.RemoveInstructorFromOfferingUseCase
	logger                     *logger.Logger
}

func NewCourseOfferingHandler(
	createOfferingUseCase *usecases.CreateCourseOfferingUseCase,
	updateOfferingUseCase *usecases.UpdateCourseOfferingUseCase,
	deleteOfferingUseCase *usecases.DeleteCourseOfferingUseCase,
	findOfferingUseCase *usecases.FindCourseOfferingUseCase,
	getOfferingUseCase *usecases.GetCourseOfferingUseCase,
	assignInstructorUseCase *usecases.AssignInstructorToOfferingUseCase,
	removeInstructorUseCase *usecases.RemoveInstructorFromOfferingUseCase,
	logger *logger.Logger,
) *CourseOfferingHandler {
	return &CourseOfferingHandler{
		createOfferingUseCase:    createOfferingUseCase,
		updateOfferingUseCase:    updateOfferingUseCase,
		deleteOfferingUseCase:    deleteOfferingUseCase,
		findOfferingUseCase:      findOfferingUseCase,
		getOfferingUseCase:       getOfferingUseCase,
		assignInstructorUseCase:  assignInstructorUseCase,
		removeInstructorUseCase: removeInstructorUseCase,
		logger:                   logger,
	}
}

// CreateCourseOffering godoc
// @Summary Create a new course offering
// @Description Create a new course offering for a course
// @Tags course-offerings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path string true "Course ID"
// @Param offering body dtos.CreateCourseOfferingInput true "Course offering creation data"
// @Success 201 {object} dtos.CourseOfferingDTO "Course offering created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{course_id}/offerings [post]
func (h *CourseOfferingHandler) CreateCourseOffering(c *gin.Context) {
	courseID := c.Param("course_id")
	if courseID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "course id is required")
		return
	}
	if _, err := uuid.Parse(courseID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course not found")
		return
	}

	var input dtos.CreateCourseOfferingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.createOfferingUseCase.Execute(c.Request.Context(), courseID, input)
	if err != nil {
		if err == usecases.ErrCourseNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to create course offering: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

// UpdateCourseOffering godoc
// @Summary Update a course offering
// @Description Update an existing course offering
// @Tags course-offerings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path string true "Course Offering ID"
// @Param offering body dtos.UpdateCourseOfferingInput true "Course offering update data"
// @Success 200 {object} dtos.CourseOfferingDTO "Course offering updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Course offering not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings/{offering_id} [put]
func (h *CourseOfferingHandler) UpdateCourseOffering(c *gin.Context) {
	offeringID := c.Param("offering_id")
	if offeringID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "offering id is required")
		return
	}
	if _, err := uuid.Parse(offeringID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course offering not found")
		return
	}

	var input dtos.UpdateCourseOfferingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.updateOfferingUseCase.Execute(c.Request.Context(), offeringID, input)
	if err != nil {
		if err == usecases.ErrCourseOfferingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to update course offering: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// FindCourseOffering godoc
// @Summary Find course offerings with filters
// @Description Find course offerings with optional filters, pagination, and sorting. Requires student, instructor, or admin role.
// @Tags course-offerings
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search query" example:"Introduction to Programming"
// @Param course_id query string false "Filter by course ID" Format(uuid)
// @Param limit query int false "Number of results per page" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of results to skip" default(0) minimum(0)
// @Param sort_column query string false "Column to sort by" Enums(name, created_at, updated_at)
// @Param sort_direction query string false "Sort direction" Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Course offerings found successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings [get]
func (h *CourseOfferingHandler) FindCourseOffering(c *gin.Context) {
	var input usecases.FindCourseOfferingInput

	if searchQuery := c.Query("search"); searchQuery != "" {
		input.SearchQuery = &searchQuery
	}

	if courseID := c.Query("course_id"); courseID != "" {
		input.CourseID = &courseID
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			input.Limit = &limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			input.Offset = &offset
		}
	}

	if sortColumn := c.Query("sort_column"); sortColumn != "" {
		input.SortColumn = &sortColumn
	}

	if sortDirection := c.Query("sort_direction"); sortDirection != "" {
		direction := repositories.SortDirection(sortDirection)
		input.SortDirection = &direction
	}

	result, err := h.findOfferingUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to find course offerings: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCourseOffering godoc
// @Summary Get course offering by ID
// @Description Retrieve course offering information with sections, modules, and instructors
// @Tags course-offerings
// @Produce json
// @Security BearerAuth
// @Param offering_id path string true "Course Offering ID" Format(uuid)
// @Success 200 {object} dtos.CourseOfferingDetailDTO "Course offering retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid offering ID"
// @Failure 404 {object} map[string]interface{} "Course offering not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings/{offering_id} [get]
func (h *CourseOfferingHandler) GetCourseOffering(c *gin.Context) {
	offeringID := c.Param("offering_id")
	if offeringID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "offering id is required")
		return
	}
	if _, err := uuid.Parse(offeringID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course offering not found")
		return
	}

	result, err := h.getOfferingUseCase.Execute(c.Request.Context(), offeringID)
	if err != nil {
		if err == usecases.ErrCourseOfferingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get course offering: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// AssignInstructor godoc
// @Summary Assign an instructor to a course offering
// @Description Assign an instructor to a course offering
// @Tags course-offerings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path string true "Course Offering ID"
// @Param instructor body dtos.AssignInstructorInput true "Instructor assignment data"
// @Success 201 {object} dtos.CourseOfferingInstructorDTO "Instructor assigned successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Course offering not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings/{offering_id}/instructors [post]
func (h *CourseOfferingHandler) AssignInstructor(c *gin.Context) {
	offeringID := c.Param("offering_id")
	if offeringID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "offering id is required")
		return
	}
	if _, err := uuid.Parse(offeringID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course offering not found")
		return
	}

	var input dtos.AssignInstructorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.assignInstructorUseCase.Execute(c.Request.Context(), offeringID, input, input.InstructorUsername)
	if err != nil {
		if err == usecases.ErrCourseOfferingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to assign instructor: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

// RemoveInstructor godoc
// @Summary Remove an instructor from a course offering
// @Description Remove an instructor from a course offering
// @Tags course-offerings
// @Produce json
// @Security BearerAuth
// @Param offering_id path string true "Course Offering ID"
// @Param instructor_id path string true "Instructor ID"
// @Success 200 {object} map[string]interface{} "Instructor removed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request parameters"
// @Failure 404 {object} map[string]interface{} "Instructor or course offering not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings/{offering_id}/instructors/{instructor_id} [delete]
func (h *CourseOfferingHandler) RemoveInstructor(c *gin.Context) {
	offeringID := c.Param("offering_id")
	if offeringID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "offering id is required")
		return
	}
	if _, err := uuid.Parse(offeringID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course offering not found")
		return
	}

	instructorID := c.Param("instructor_id")
	if instructorID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "instructor id is required")
		return
	}
	if _, err := uuid.Parse(instructorID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "instructor not found")
		return
	}

	err := h.removeInstructorUseCase.Execute(c.Request.Context(), offeringID, instructorID)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to remove instructor: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "instructor removed successfully"})
}

// DeleteCourseOffering godoc
// @Summary Delete a course offering
// @Description Delete an existing course offering
// @Tags course-offerings
// @Produce json
// @Security BearerAuth
// @Param offering_id path string true "Course Offering ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Course offering deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid offering ID"
// @Failure 404 {object} map[string]interface{} "Course offering not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings/{offering_id} [delete]
func (h *CourseOfferingHandler) DeleteCourseOffering(c *gin.Context) {
	offeringID := c.Param("offering_id")
	if offeringID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "offering id is required")
		return
	}
	if _, err := uuid.Parse(offeringID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course offering not found")
		return
	}

	err := h.deleteOfferingUseCase.Execute(c.Request.Context(), offeringID)
	if err != nil {
		if err == usecases.ErrCourseOfferingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to delete course offering: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "course offering deleted successfully"})
}

