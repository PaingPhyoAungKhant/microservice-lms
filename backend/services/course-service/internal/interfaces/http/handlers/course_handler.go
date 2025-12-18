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

type CourseHandler struct {
	createCourseUseCase      *usecases.CreateCourseUseCase
	listCoursesUseCase       *usecases.ListCoursesUseCase
	findCourseUseCase        *usecases.FindCourseUseCase
	getCourseUseCase         *usecases.GetCourseUseCase
	getCourseWithDetailsUseCase *usecases.GetCourseWithDetailsUseCase
	updateCourseUseCase      *usecases.UpdateCourseUseCase
	deleteCourseUseCase      *usecases.DeleteCourseUseCase
	logger                   *logger.Logger
}

func NewCourseHandler(
	createCourseUseCase *usecases.CreateCourseUseCase,
	listCoursesUseCase *usecases.ListCoursesUseCase,
	findCourseUseCase *usecases.FindCourseUseCase,
	getCourseUseCase *usecases.GetCourseUseCase,
	getCourseWithDetailsUseCase *usecases.GetCourseWithDetailsUseCase,
	updateCourseUseCase *usecases.UpdateCourseUseCase,
	deleteCourseUseCase *usecases.DeleteCourseUseCase,
	logger *logger.Logger,
) *CourseHandler {
	return &CourseHandler{
		createCourseUseCase:      createCourseUseCase,
		listCoursesUseCase:       listCoursesUseCase,
		findCourseUseCase:        findCourseUseCase,
		getCourseUseCase:         getCourseUseCase,
		getCourseWithDetailsUseCase: getCourseWithDetailsUseCase,
		updateCourseUseCase:      updateCourseUseCase,
		deleteCourseUseCase:      deleteCourseUseCase,
		logger:                   logger,
	}
}

// CreateCourse godoc
// @Summary Create a new course
// @Description Create a new course. Requires admin role.
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course body dtos.CreateCourseInput true "Course creation data"
// @Success 201 {object} dtos.CourseDTO "Course created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var input dtos.CreateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.createCourseUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to create course: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *CourseHandler) ListCourses(c *gin.Context) {
	var input usecases.ListCoursesInput

	if searchQuery := c.Query("search"); searchQuery != "" {
		input.SearchQuery = &searchQuery
	}

	if categoryID := c.Query("category_id"); categoryID != "" {
		input.CategoryID = &categoryID
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

	result, err := h.listCoursesUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to list courses: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": result.Courses,
		"total":   result.Total,
	})
}

// FindCourse godoc
// @Summary Find courses with filters
// @Description Find courses with optional filters, pagination, and sorting. Requires student, instructor, or admin role.
// @Tags courses
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search query" example:"Introduction to Programming"
// @Param category_id query string false "Filter by category ID" Format(uuid)
// @Param limit query int false "Number of results per page" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of results to skip" default(0) minimum(0)
// @Param sort_column query string false "Column to sort by" Enums(name, created_at, updated_at)
// @Param sort_direction query string false "Sort direction" Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Courses found successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses [get]
func (h *CourseHandler) FindCourse(c *gin.Context) {
	var input usecases.FindCourseInput

	if searchQuery := c.Query("search"); searchQuery != "" {
		input.SearchQuery = &searchQuery
	}

	if categoryID := c.Query("category_id"); categoryID != "" {
		input.CategoryID = &categoryID
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

	result, err := h.findCourseUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to find courses: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCourse godoc
// @Summary Get course by ID
// @Description Retrieve course information by course ID. Requires student, instructor, or admin role.
// @Tags courses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Course ID" Format(uuid) example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} dtos.CourseDTO "Course retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "course id is required")
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course not found")
		return
	}

	input := usecases.GetCourseInput{
		CourseID: id,
	}

	result, err := h.getCourseUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrCourseNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get course: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCourseWithDetails godoc
// @Summary Get course with all details
// @Description Retrieve course information with all offerings, sections, modules, and instructors
// @Tags courses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Course ID" Format(uuid) example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} dtos.CourseDetailDTO "Course details retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/details [get]
func (h *CourseHandler) GetCourseWithDetails(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "course id is required")
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course not found")
		return
	}

	result, err := h.getCourseWithDetailsUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrCourseNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get course details: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateCourse godoc
// @Summary Update a course by ID
// @Description Update course information by course ID. Requires admin role.
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Course ID" Format(uuid) example:"550e8400-e29b-41d4-a716-446655440000"
// @Param course body dtos.UpdateCourseInput true "Course update data"
// @Success 200 {object} dtos.CourseDTO "Course updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or course ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "course id is required")
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course not found")
		return
	}

	var input dtos.UpdateCourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	updateInput := usecases.UpdateCourseInput{
		CourseID: id,
	}

	result, err := h.updateCourseUseCase.Execute(c.Request.Context(), updateInput, input)
	if err != nil {
		if err == usecases.ErrCourseNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to update course: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteCourse godoc
// @Summary Delete a course by ID
// @Description Delete a course by its ID. Requires admin role.
// @Tags courses
// @Produce json
// @Security BearerAuth
// @Param id path string true "Course ID" Format(uuid) example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} map[string]interface{} "Course deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "course id is required")
		return
	}

	input := usecases.DeleteCourseInput{
		CourseID: id,
	}

	result, err := h.deleteCourseUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrCourseNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to delete course: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

