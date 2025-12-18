package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/user-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/domain/valueobjects"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
)

type UserHandler struct {
	createUserUseCase *usecases.CreateUserUseCase
	getUserUseCase *usecases.GetUserUseCase
	updateUserUseCase *usecases.UpdateUserUseCase
	findUserUseCase *usecases.FindUserUseCase
	deleteUserUseCase *usecases.DeleteUserUsecase
	logger *logger.Logger
}

func NewUserHandler(
	createUserUseCase *usecases.CreateUserUseCase, 
	getUserUseCase *usecases.GetUserUseCase, 
	updateUserUseCase *usecases.UpdateUserUseCase, 
	findUserUseCase *usecases.FindUserUseCase,
	deleteUserUseCase *usecases.DeleteUserUsecase,
	logger *logger.Logger,
	) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
		getUserUseCase: getUserUseCase,
		updateUserUseCase: updateUserUseCase,
		findUserUseCase: findUserUseCase,
		deleteUserUseCase: deleteUserUseCase,
		logger: logger,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"paingkhant0397@gmail.com"`
	Username string `json:"username" binding:"required,min=3,max=255" example:"testuser"`
	Password string `json:"password" binding:"required,min=8,max=255" example:"Password@123"`
	Role     string `json:"role" binding:"required,oneof=student instructor admin" example:"student"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account. Requires admin role.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User creation details"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required"
// @Failure 409 {object} map[string]interface{} "Email or username already exists"
// @Router / [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest 
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	input := usecases.CreateUserInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	output, err := h.createUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrUsernameAlreadyExists || err == usecases.ErrEmailAlreadyExists {
			middleware.AbortWithError(c, http.StatusConflict, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, output)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve user information by user ID. Requires admin or instructor role.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "User retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin or instructor role required"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	input := usecases.GetUserInput{
		UserID: id,
	}

	output, err := h.getUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrUserNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=255" example:"updateduser"`
	Email    string `json:"email" binding:"omitempty,email" example:"updated@example.com"`
	Role     string `json:"role" binding:"omitempty,oneof=student instructor admin" example:"student"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive pending banned" example:"active"`
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user information. Requires admin role.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" Format(uuid)
// @Param request body UpdateUserRequest true "User update details"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	input := usecases.UpdateUserInput{
		ID: id,
		Username: &req.Username,
		Email: &req.Email,
		Role: &req.Role,
		Status: &req.Status,
	}

	output, err := h.updateUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrUserNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, output)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user account. Requires admin role. Cannot delete admin users.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin role required or cannot delete admin user"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	input := usecases.DeleteUserInput{
		UserID: id,
	}

	output, err := h.deleteUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrUserNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, output)
}

type FindUserRequest struct {
	SearchQuery   string `json:"search_query" form:"search_query" binding:"omitempty,min=1,max=255" example:"testuser"`
	Role          string `json:"role" form:"role" binding:"omitempty,oneof=student instructor admin" example:"student"`
	Status        string `json:"status" form:"status" binding:"omitempty,oneof=active inactive pending banned" example:"active"`
	Limit         int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=100" example:"10"`
	Offset        int    `json:"offset" form:"offset" binding:"omitempty,min=0" example:"0"`
	SortColumn    string `json:"sort_column" form:"sort_column" binding:"omitempty,oneof=username email role status created_at updated_at" example:"created_at"`
	SortDirection string `json:"sort_direction" form:"sort_direction" binding:"omitempty,oneof=asc desc" example:"desc"`
}

// FindUser godoc
// @Summary Find users with filters
// @Description Search and filter users with pagination and sorting. Requires admin or instructor role.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search_query query string false "Search query for username or email"
// @Param role query string false "Filter by role" Enums(student, instructor, admin)
// @Param status query string false "Filter by status" Enums(active, inactive, pending, banned)
// @Param limit query int false "Number of results per page" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of results to skip" default(0) minimum(0)
// @Param sort_column query string false "Column to sort by" Enums(username, email, role, status, created_at, updated_at)
// @Param sort_direction query string false "Sort direction" Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Users retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Admin or instructor role required"
// @Router / [get]
func (h *UserHandler) FindUser(c *gin.Context) { 
	var req FindUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	input := usecases.FindUserInput{}

	if req.SearchQuery != "" {
		input.SearchQuery = &req.SearchQuery
	}

	if req.Role != "" {
		role := valueobjects.Role(req.Role)
		input.Role = &role
	}

	if req.Status != "" {
		status := valueobjects.Status(req.Status)
		input.Status = &status
	}

	if req.Limit > 0 {
		input.Limit = &req.Limit
	}

	// Set offset when limit is provided (pagination) or when offset > 0
	// When limit is set, always set offset (even if 0 for first page)
	// Note: The repository will only add OFFSET clause if offset > 0
	if req.Limit > 0 {
		// When limit is set, always set offset (even if 0 for first page)
		input.Offset = &req.Offset
	} else if req.Offset > 0 {
		// If limit not set but offset > 0, set it
		input.Offset = &req.Offset
	}

	if req.SortColumn != "" {
		input.SortColumn = &req.SortColumn
	}

	if req.SortDirection != "" {
		// Normalize to uppercase to match repository constants
		sortDir := repositories.SortDirection(strings.ToUpper(req.SortDirection))
		input.SortDirection = &sortDir
	}

	output, err := h.findUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

// Health godoc
// @Summary Health check
// @Description Check if the user service is running and healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Service is healthy"
// @Router /health [get]
func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "user-service",
	})
}