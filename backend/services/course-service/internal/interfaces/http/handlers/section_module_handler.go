package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/application/dtos"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
)

type SectionModuleHandler struct {
	createModuleUseCase *usecases.CreateSectionModuleUseCase
	updateModuleUseCase *usecases.UpdateSectionModuleUseCase
	getModuleUseCase    *usecases.GetSectionModuleUseCase
	deleteModuleUseCase *usecases.DeleteSectionModuleUseCase
	findModuleUseCase   *usecases.FindSectionModuleUseCase
	reorderModulesUseCase *usecases.ReorderSectionModulesUseCase
	logger              *logger.Logger
}

func NewSectionModuleHandler(
	createModuleUseCase *usecases.CreateSectionModuleUseCase,
	updateModuleUseCase *usecases.UpdateSectionModuleUseCase,
	getModuleUseCase *usecases.GetSectionModuleUseCase,
	deleteModuleUseCase *usecases.DeleteSectionModuleUseCase,
	findModuleUseCase *usecases.FindSectionModuleUseCase,
	reorderModulesUseCase *usecases.ReorderSectionModulesUseCase,
	logger *logger.Logger,
) *SectionModuleHandler {
	return &SectionModuleHandler{
		createModuleUseCase: createModuleUseCase,
		updateModuleUseCase: updateModuleUseCase,
		getModuleUseCase:    getModuleUseCase,
		deleteModuleUseCase: deleteModuleUseCase,
		findModuleUseCase:   findModuleUseCase,
		reorderModulesUseCase: reorderModulesUseCase,
		logger:              logger,
	}
}

// CreateSectionModule godoc
// @Summary Create a new section module
// @Description Create a new module for a course section
// @Tags section-modules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param section_id path string true "Course Section ID"
// @Param module body dtos.CreateSectionModuleInput true "Section module creation data"
// @Success 201 {object} dtos.SectionModuleDTO "Section module created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Course section not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-sections/{section_id}/modules [post]
func (h *SectionModuleHandler) CreateSectionModule(c *gin.Context) {
	sectionID := c.Param("section_id")
	if sectionID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "section id is required")
		return
	}
	if _, err := uuid.Parse(sectionID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course section not found")
		return
	}

	var input dtos.CreateSectionModuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.createModuleUseCase.Execute(c.Request.Context(), sectionID, input)
	if err != nil {
		if err == usecases.ErrCourseSectionNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to create section module: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

// UpdateSectionModule godoc
// @Summary Update a section module
// @Description Update an existing section module
// @Tags section-modules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param module_id path string true "Section Module ID"
// @Param module body dtos.UpdateSectionModuleInput true "Section module update data"
// @Success 200 {object} dtos.SectionModuleDTO "Section module updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Section module not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /section-modules/{module_id} [put]
func (h *SectionModuleHandler) UpdateSectionModule(c *gin.Context) {
	moduleID := c.Param("module_id")
	if moduleID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "module id is required")
		return
	}
	if _, err := uuid.Parse(moduleID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "section module not found")
		return
	}

	var input dtos.UpdateSectionModuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.updateModuleUseCase.Execute(c.Request.Context(), moduleID, input)
	if err != nil {
		if err == usecases.ErrSectionModuleNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to update section module: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSectionModule godoc
// @Summary Get section module by ID
// @Description Retrieve section module information by module ID
// @Tags section-modules
// @Produce json
// @Security BearerAuth
// @Param module_id path string true "Section Module ID" Format(uuid)
// @Success 200 {object} dtos.SectionModuleDTO "Section module retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid module ID"
// @Failure 404 {object} map[string]interface{} "Section module not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /section-modules/{module_id} [get]
func (h *SectionModuleHandler) GetSectionModule(c *gin.Context) {
	moduleID := c.Param("module_id")
	if moduleID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "module id is required")
		return
	}
	if _, err := uuid.Parse(moduleID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "section module not found")
		return
	}

	result, err := h.getModuleUseCase.Execute(c.Request.Context(), usecases.GetSectionModuleInput{ModuleID: moduleID})
	if err != nil {
		if err == usecases.ErrSectionModuleNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get section module: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteSectionModule godoc
// @Summary Delete a section module
// @Description Delete an existing section module
// @Tags section-modules
// @Produce json
// @Security BearerAuth
// @Param module_id path string true "Section Module ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Section module deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid module ID"
// @Failure 404 {object} map[string]interface{} "Section module not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /section-modules/{module_id} [delete]
func (h *SectionModuleHandler) DeleteSectionModule(c *gin.Context) {
	moduleID := c.Param("module_id")
	if moduleID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "module id is required")
		return
	}
	if _, err := uuid.Parse(moduleID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "section module not found")
		return
	}

	result, err := h.deleteModuleUseCase.Execute(c.Request.Context(), usecases.DeleteSectionModuleInput{ModuleID: moduleID})
	if err != nil {
		if err == usecases.ErrSectionModuleNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to delete section module: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// FindSectionModule godoc
// @Summary Find section modules by section ID
// @Description Retrieve all section modules for a specific course section
// @Tags section-modules
// @Produce json
// @Security BearerAuth
// @Param section_id path string true "Course Section ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Section modules retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid section ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-sections/{section_id}/modules [get]
func (h *SectionModuleHandler) FindSectionModule(c *gin.Context) {
	sectionID := c.Param("section_id")
	if sectionID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "section id is required")
		return
	}
	if _, err := uuid.Parse(sectionID); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "invalid section id")
		return
	}

	result, err := h.findModuleUseCase.Execute(c.Request.Context(), usecases.FindSectionModuleInput{SectionID: sectionID})
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to find section modules: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// ReorderSectionModules godoc
// @Summary Reorder section modules
// @Description Reorder section modules for a course section
// @Tags section-modules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param section_id path string true "Course Section ID" Format(uuid)
// @Param items body usecases.ReorderSectionModulesInput true "Reorder items with module_id and order"
// @Success 200 {object} map[string]interface{} "Section modules reordered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Course section or module not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-sections/{section_id}/modules/reorder [put]
func (h *SectionModuleHandler) ReorderSectionModules(c *gin.Context) {
	sectionID := c.Param("section_id")
	if sectionID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "section id is required")
		return
	}
	if _, err := uuid.Parse(sectionID); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "invalid section id")
		return
	}

	var input usecases.ReorderSectionModulesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	input.SectionID = sectionID

	err := h.reorderModulesUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrSectionModuleNotFound || err == usecases.ErrInvalidReorderInput {
			middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to reorder section modules: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "section modules reordered successfully"})
}

