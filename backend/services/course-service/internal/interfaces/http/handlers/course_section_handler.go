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

type CourseSectionHandler struct {
	createSectionUseCase *usecases.CreateCourseSectionUseCase
	updateSectionUseCase *usecases.UpdateCourseSectionUseCase
	getSectionUseCase    *usecases.GetCourseSectionUseCase
	deleteSectionUseCase *usecases.DeleteCourseSectionUseCase
	findSectionUseCase   *usecases.FindCourseSectionUseCase
	reorderSectionsUseCase *usecases.ReorderCourseSectionsUseCase
	logger               *logger.Logger
}

func NewCourseSectionHandler(
	createSectionUseCase *usecases.CreateCourseSectionUseCase,
	updateSectionUseCase *usecases.UpdateCourseSectionUseCase,
	getSectionUseCase *usecases.GetCourseSectionUseCase,
	deleteSectionUseCase *usecases.DeleteCourseSectionUseCase,
	findSectionUseCase *usecases.FindCourseSectionUseCase,
	reorderSectionsUseCase *usecases.ReorderCourseSectionsUseCase,
	logger *logger.Logger,
) *CourseSectionHandler {
	return &CourseSectionHandler{
		createSectionUseCase: createSectionUseCase,
		updateSectionUseCase: updateSectionUseCase,
		getSectionUseCase:    getSectionUseCase,
		deleteSectionUseCase: deleteSectionUseCase,
		findSectionUseCase:   findSectionUseCase,
		reorderSectionsUseCase: reorderSectionsUseCase,
		logger:               logger,
	}
}

// CreateCourseSection godoc
// @Summary Create a new course section
// @Description Create a new course section for a course offering
// @Tags course-sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path string true "Course Offering ID"
// @Param section body dtos.CreateCourseSectionInput true "Course section creation data"
// @Success 201 {object} dtos.CourseSectionDTO "Course section created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Course offering not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings/{offering_id}/sections [post]
func (h *CourseSectionHandler) CreateCourseSection(c *gin.Context) {
	offeringID := c.Param("offering_id")
	if offeringID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "offering id is required")
		return
	}
	if _, err := uuid.Parse(offeringID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course offering not found")
		return
	}

	var input dtos.CreateCourseSectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.createSectionUseCase.Execute(c.Request.Context(), offeringID, input)
	if err != nil {
		if err == usecases.ErrCourseOfferingNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to create course section: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, result)
}

// UpdateCourseSection godoc
// @Summary Update a course section
// @Description Update an existing course section
// @Tags course-sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param section_id path string true "Course Section ID"
// @Param section body dtos.UpdateCourseSectionInput true "Course section update data"
// @Success 200 {object} dtos.CourseSectionDTO "Course section updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Course section not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-sections/{section_id} [put]
func (h *CourseSectionHandler) UpdateCourseSection(c *gin.Context) {
	sectionID := c.Param("section_id")
	if sectionID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "section id is required")
		return
	}
	if _, err := uuid.Parse(sectionID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course section not found")
		return
	}

	var input dtos.UpdateCourseSectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.updateSectionUseCase.Execute(c.Request.Context(), sectionID, input)
	if err != nil {
		if err == usecases.ErrCourseSectionNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to update course section: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCourseSection godoc
// @Summary Get course section by ID
// @Description Retrieve course section information by section ID
// @Tags course-sections
// @Produce json
// @Security BearerAuth
// @Param section_id path string true "Course Section ID" Format(uuid)
// @Success 200 {object} dtos.CourseSectionDTO "Course section retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid section ID"
// @Failure 404 {object} map[string]interface{} "Course section not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-sections/{section_id} [get]
func (h *CourseSectionHandler) GetCourseSection(c *gin.Context) {
	sectionID := c.Param("section_id")
	if sectionID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "section id is required")
		return
	}
	if _, err := uuid.Parse(sectionID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course section not found")
		return
	}

	result, err := h.getSectionUseCase.Execute(c.Request.Context(), usecases.GetCourseSectionInput{SectionID: sectionID})
	if err != nil {
		if err == usecases.ErrCourseSectionNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get course section: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteCourseSection godoc
// @Summary Delete a course section
// @Description Delete an existing course section
// @Tags course-sections
// @Produce json
// @Security BearerAuth
// @Param section_id path string true "Course Section ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Course section deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid section ID"
// @Failure 404 {object} map[string]interface{} "Course section not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-sections/{section_id} [delete]
func (h *CourseSectionHandler) DeleteCourseSection(c *gin.Context) {
	sectionID := c.Param("section_id")
	if sectionID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "section id is required")
		return
	}
	if _, err := uuid.Parse(sectionID); err != nil {
		middleware.AbortWithError(c, http.StatusNotFound, "course section not found")
		return
	}

	result, err := h.deleteSectionUseCase.Execute(c.Request.Context(), usecases.DeleteCourseSectionInput{SectionID: sectionID})
	if err != nil {
		if err == usecases.ErrCourseSectionNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to delete course section: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// FindCourseSection godoc
// @Summary Find course sections by offering ID
// @Description Retrieve all course sections for a specific course offering
// @Tags course-sections
// @Produce json
// @Security BearerAuth
// @Param offering_id path string true "Course Offering ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "Course sections retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid offering ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings/{offering_id}/sections [get]
func (h *CourseSectionHandler) FindCourseSection(c *gin.Context) {
	offeringID := c.Param("offering_id")
	if offeringID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "offering id is required")
		return
	}
	if _, err := uuid.Parse(offeringID); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "invalid offering id")
		return
	}

	result, err := h.findSectionUseCase.Execute(c.Request.Context(), usecases.FindCourseSectionInput{OfferingID: offeringID})
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to find course sections: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// ReorderCourseSections godoc
// @Summary Reorder course sections
// @Description Reorder course sections for a course offering
// @Tags course-sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path string true "Course Offering ID" Format(uuid)
// @Param items body usecases.ReorderCourseSectionsInput true "Reorder items with section_id and order"
// @Success 200 {object} map[string]interface{} "Course sections reordered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 404 {object} map[string]interface{} "Course offering or section not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /course-offerings/{offering_id}/sections/reorder [put]
func (h *CourseSectionHandler) ReorderCourseSections(c *gin.Context) {
	offeringID := c.Param("offering_id")
	if offeringID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "offering id is required")
		return
	}
	if _, err := uuid.Parse(offeringID); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "invalid offering id")
		return
	}

	var input usecases.ReorderCourseSectionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	input.OfferingID = offeringID

	err := h.reorderSectionsUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrCourseSectionNotFound || err == usecases.ErrInvalidReorderInput {
			middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to reorder course sections: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "course sections reordered successfully"})
}

