package dtos

import (
	"time"

	"github.com/paingphyoaungkhant/asto-microservice/services/file-service/internal/domain/entities"
)

type FileDTO struct {
	ID               string    `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	StoredFilename   string    `json:"stored_filename"`
	BucketName       string    `json:"bucket_name"`
	MimeType         string    `json:"mime_type"`
	SizeBytes        int64     `json:"size_bytes"`
	UploadedBy       string    `json:"uploaded_by"`
	Tags             []string  `json:"tags"`
	DownloadURL      string    `json:"download_url,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

func (d *FileDTO) FromEntity(file *entities.File, apiGatewayURL string) {
	d.ID = file.ID
	d.OriginalFilename = file.OriginalFilename
	d.StoredFilename = file.StoredFilename
	d.BucketName = file.BucketName
	d.MimeType = file.MimeType
	d.SizeBytes = file.SizeBytes
	d.UploadedBy = file.UploadedBy
	d.Tags = file.Tags
	d.CreatedAt = file.CreatedAt
	d.UpdatedAt = file.UpdatedAt
	d.DeletedAt = file.DeletedAt
	
	if apiGatewayURL != "" {
		d.DownloadURL = apiGatewayURL + "/api/v1/files/" + file.ID + "/download"
	}
}

