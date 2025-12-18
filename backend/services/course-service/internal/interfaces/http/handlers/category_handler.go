package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
)

type CategoryHandler struct {
	createCategoryUseCase *usecases.CreateCategoryUseCase
	listCategoriesUseCase *usecases.ListCategoriesUseCase
	findCategoryUseCase   *usecases.FindCategoryUseCase
	getCategoryUseCase    *usecases.GetCategoryUseCase
	updateCategoryUseCase *usecases.UpdateCategoryUseCase
	deleteCategoryUseCase *usecases.DeleteCategoryUseCase
	logger                *logger.Logger
}

func NewCategoryHandler(
	createCategoryUseCase *usecases.CreateCategoryUseCase,
	listCategoriesUseCase *usecases.ListCategoriesUseCase,
	findCategoryUseCase *usecases.FindCategoryUseCase,
	getCategoryUseCase *usecases.GetCategoryUseCase,
	updateCategoryUseCase *usecases.UpdateCategoryUseCase,
	deleteCategoryUseCase *usecases.DeleteCategoryUseCase,
	logger *logger.Logger,
) *CategoryHandler {
	return &CategoryHandler{
		createCategoryUseCase: createCategoryUseCase,
		listCategoriesUseCase: listCategoriesUseCase,
		findCategoryUseCase:   findCategoryUseCase,
		getCategoryUseCase:    getCategoryUseCase,
		updateCategoryUseCase: updateCategoryUseCase,
		deleteCategoryUseCase: deleteCategoryUseCase,
		logger:                logger,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category. Requires admin role.
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body dtos.CreateCategoryInput true "Category creation data"
// @Success 201 {object} dtos.CategoryDTO "Category created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required"
// @Failure 409 {object} map[string]interface{} "Category already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input dtos.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.createCategoryUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrCategoryAlreadyExists {
			middleware.AbortWithError(c, http.StatusConflict, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to create category: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	var input usecases.ListCategoriesInput

	if searchQuery := c.Query("search"); searchQuery != "" {
		input.SearchQuery = &searchQuery
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

	result, err := h.listCategoriesUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to list categories: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": result.Categories,
		"total":      result.Total,
	})
}

// FindCategory godoc
// @Summary Find categories with filters
// @Description Find categories with optional filters, pagination, and sorting. Requires student, instructor, or admin role.
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search query" example:"Programming"
// @Param limit query int false "Number of results per page" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of results to skip" default(0) minimum(0)
// @Param sort_column query string false "Column to sort by" Enums(name, created_at, updated_at)
// @Param sort_direction query string false "Sort direction" Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Categories found successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) FindCategory(c *gin.Context) {
	var input usecases.FindCategoryInput

	if searchQuery := c.Query("search"); searchQuery != "" {
		input.SearchQuery = &searchQuery
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

	result, err := h.findCategoryUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to find categories: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCategory godoc
// @Summary Get category by ID
// @Description Retrieve category information by category ID. Requires student, instructor, or admin role.
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID" Format(uuid) example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} dtos.CategoryDTO "Category retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid category ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "category id is required")
		return
	}

	input := usecases.GetCategoryInput{
		CategoryID: id,
	}

	result, err := h.getCategoryUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrCategoryNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get category: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateCategory godoc
// @Summary Update a category by ID
// @Description Update category information by category ID. Requires admin role.
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID" Format(uuid) example:"550e8400-e29b-41d4-a716-446655440000"
// @Param category body dtos.UpdateCategoryInput true "Category update data"
// @Success 200 {object} dtos.CategoryDTO "Category updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or category ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "category id is required")
		return
	}

	var input dtos.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	updateInput := usecases.UpdateCategoryInput{
		CategoryID: id,
	}

	result, err := h.updateCategoryUseCase.Execute(c.Request.Context(), updateInput, input)
	if err != nil {
		if err == usecases.ErrCategoryNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to update category: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteCategory godoc
// @Summary Delete a category by ID
// @Description Delete a category by its ID. Requires admin role.
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID" Format(uuid) example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid category ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "category id is required")
		return
	}

	input := usecases.DeleteCategoryInput{
		CategoryID: id,
	}

	result, err := h.deleteCategoryUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrCategoryNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to delete category: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

