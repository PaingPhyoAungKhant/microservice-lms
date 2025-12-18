package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/application/usecases"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"github.com/paingphyoaungkhant/asto-microservice/shared/middleware"
	"go.uber.org/zap"
)

type FileHandler struct {
	uploadFileUseCase   *usecases.UploadFileUseCase
	downloadFileUseCase *usecases.DownloadFileUseCase
	getFileUseCase      *usecases.GetFileUseCase
	listFilesUseCase    *usecases.ListFilesUseCase
	deleteFileUseCase   *usecases.DeleteFileUseCase
	logger              *logger.Logger
}

func NewFileHandler(
	uploadFileUseCase *usecases.UploadFileUseCase,
	downloadFileUseCase *usecases.DownloadFileUseCase,
	getFileUseCase *usecases.GetFileUseCase,
	listFilesUseCase *usecases.ListFilesUseCase,
	deleteFileUseCase *usecases.DeleteFileUseCase,
	logger *logger.Logger,
) *FileHandler {
	return &FileHandler{
		uploadFileUseCase:   uploadFileUseCase,
		downloadFileUseCase: downloadFileUseCase,
		getFileUseCase:      getFileUseCase,
		listFilesUseCase:    listFilesUseCase,
		deleteFileUseCase:   deleteFileUseCase,
		logger:              logger,
	}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a new file. Requires instructor or admin role.
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param bucket_name formData string false "Bucket name" example:"course-thumbnails"
// @Param tags formData string false "Comma-separated tags" example:"course,thumbnail"
// @Success 201 {object} map[string]interface{} "File uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Instructor or admin role required"
// @Failure 413 {object} map[string]interface{} "File too large"
// @Router /files [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		middleware.AbortWithError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := c.Request.ParseMultipartForm(10 << 20) // 10MB max memory
	if err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "failed to parse multipart form: "+err.Error())
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "file is required: "+err.Error())
		return
	}
	defer file.Close()

	bucketName := c.PostForm("bucket_name")
	tagsStr := c.PostForm("tags")
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	input := usecases.UploadFileInput{
		File:       file,
		Filename:   header.Filename,
		MimeType:   contentType,
		Size:       header.Size,
		UploadedBy: userID,
		BucketName: bucketName,
		Tags:       tags,
	}

	output, err := h.uploadFileUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrFileTooLarge {
			middleware.AbortWithError(c, http.StatusRequestEntityTooLarge, err.Error())
			return
		}
		if err == usecases.ErrInvalidMimeType || err == usecases.ErrFileRequired {
			middleware.AbortWithError(c, http.StatusBadRequest, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to upload file: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, output)
}

// DownloadFile godoc
// @Summary Download a file by ID
// @Description Download a file by its ID. Requires authentication.
// @Tags files
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "File ID" Format(uuid)
// @Success 200 {file} file "File content"
// @Failure 400 {object} map[string]interface{} "Invalid file ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Router /files/{id}/download [get]
func (h *FileHandler) DownloadFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "file id is required")
		return
	}

	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		middleware.AbortWithError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	input := usecases.DownloadFileInput{
		FileID: fileID,
		UserID: userID,
	}

	output, err := h.downloadFileUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrFileNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to download file: "+err.Error())
		return
	}
	defer output.Reader.Close()

	c.Header("Content-Type", output.ContentType)
	c.Header("Content-Disposition", `attachment; filename="`+output.Filename+`"`)
	c.Header("Content-Length", strconv.FormatInt(output.ContentLength, 10))

	_, err = io.Copy(c.Writer, output.Reader)
	if err != nil {
		h.logger.Error("failed to stream file", zap.Error(err))
		return
	}
}

// GetFile godoc
// @Summary Get file information by ID
// @Description Retrieve file metadata by file ID.
// @Tags files
// @Produce json
// @Param id path string true "File ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "File information"
// @Failure 400 {object} map[string]interface{} "Invalid file ID"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Router /files/{id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "file id is required")
		return
	}

	input := usecases.GetFileInput{
		FileID: fileID,
	}

	output, err := h.getFileUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrFileNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to get file: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}

// ListFiles godoc
// @Summary List files with filters
// @Description List files with optional filters, pagination, and sorting. Requires authentication.
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param uploaded_by query string false "Filter by uploader user ID" Format(uuid)
// @Param tags query string false "Comma-separated tags to filter" example:"course,thumbnail"
// @Param mime_type query string false "Filter by MIME type" example:"image/png"
// @Param bucket_name query string false "Filter by bucket name" example:"course-thumbnails"
// @Param limit query int false "Number of results per page" default(10) minimum(1) maximum(100)
// @Param offset query int false "Number of results to skip" default(0) minimum(0)
// @Param sort_column query string false "Column to sort by" Enums(created_at, updated_at, filename, size)
// @Param sort_direction query string false "Sort direction" Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Files list"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /files [get]
func (h *FileHandler) ListFiles(c *gin.Context) {
	var input usecases.ListFilesInput

	if uploadedBy := c.Query("uploaded_by"); uploadedBy != "" {
		input.UploadedBy = &uploadedBy
	}

	if tagsStr := c.Query("tags"); tagsStr != "" {
		tags := strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		input.Tags = tags
	}

	if mimeType := c.Query("mime_type"); mimeType != "" {
		input.MimeType = &mimeType
	}

	if bucketName := c.Query("bucket_name"); bucketName != "" {
		input.BucketName = &bucketName
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
		input.SortDirection = &sortDirection
	}

	files, total, err := h.listFilesUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to list files: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
		"total": total,
	})
}

// DeleteFile godoc
// @Summary Delete a file by ID
// @Description Delete a file by its ID. Requires instructor or admin role.
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID" Format(uuid)
// @Success 200 {object} map[string]interface{} "File deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Instructor or admin role required"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Router /files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "file id is required")
		return
	}

	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		middleware.AbortWithError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	input := usecases.DeleteFileInput{
		FileID: fileID,
		UserID: userID,
	}

	err := h.deleteFileUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrFileNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to delete file: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}

// DownloadFileByBucket godoc
// @Summary Download a file from a specific bucket
// @Description Download a file by ID from a specific bucket. Course thumbnails are public (no auth required). For zoom recordings, requires student, instructor, or admin role.
// @Tags buckets
// @Produce application/octet-stream
// @Security BearerAuth
// @Param bucket path string true "Bucket name" example:"course-thumbnails"
// @Param id path string true "File ID" Format(uuid)
// @Success 200 {file} file "File content"
// @Failure 400 {object} map[string]interface{} "Invalid bucket name or file ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Router /buckets/{bucket}/files/{id}/download [get]
func (h *FileHandler) DownloadFileByBucket(c *gin.Context) {
	bucketName := c.Param("bucket")
	fileID := c.Param("id")
	
	if bucketName == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "bucket name is required")
		return
	}
	
	if fileID == "" {
		middleware.AbortWithError(c, http.StatusBadRequest, "file id is required")
		return
	}

	userID := c.GetHeader("X-User-ID")
	
	if bucketName != "course-thumbnails" && userID == "" {
		middleware.AbortWithError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	input := usecases.DownloadFileInput{
		FileID: fileID,
		UserID: userID,
		BucketName: &bucketName,
	}

	output, err := h.downloadFileUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if err == usecases.ErrFileNotFound {
			middleware.AbortWithError(c, http.StatusNotFound, err.Error())
			return
		}
		middleware.AbortWithError(c, http.StatusInternalServerError, "failed to download file: "+err.Error())
		return
	}
	defer output.Reader.Close()

	c.Header("Content-Type", output.ContentType)
	c.Header("Content-Disposition", `attachment; filename="`+output.Filename+`"`)
	c.Header("Content-Length", strconv.FormatInt(output.ContentLength, 10))

	_, err = io.Copy(c.Writer, output.Reader)
	if err != nil {
		h.logger.Error("failed to stream file", zap.Error(err))
		return
	}
}

// Health godoc
// @Summary Health check
// @Description Check the health status of the file service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "Service is healthy"
// @Router /health [get]
func (h *FileHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "file-service",
	})
}

