package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID             string
	OriginalFilename string
	StoredFilename   string
	BucketName       string
	MimeType         string
	SizeBytes        int64
	UploadedBy       string
	Tags             []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

func NewFile(originalFilename, storedFilename, bucketName, mimeType string, sizeBytes int64, uploadedBy string, tags []string) *File {
	now := time.Now().UTC()
	return &File{
		ID:               uuid.NewString(),
		OriginalFilename: originalFilename,
		StoredFilename:   storedFilename,
		BucketName:       bucketName,
		MimeType:         mimeType,
		SizeBytes:        sizeBytes,
		UploadedBy:       uploadedBy,
		Tags:             tags,
		CreatedAt:        now,
		UpdatedAt:        now,
		DeletedAt:        nil,
	}
}

func (f *File) UpdateTags(tags []string) {
	f.Tags = tags
	f.UpdatedAt = time.Now().UTC()
}

func (f *File) SoftDelete() {
	now := time.Now().UTC()
	f.DeletedAt = &now
	f.UpdatedAt = now
}

func (f *File) IsDeleted() bool {
	return f.DeletedAt != nil
}

func (f *File) TagsToJSONB() (json.RawMessage, error) {
	if f.Tags == nil {
		return json.RawMessage("[]"), nil
	}
	return json.Marshal(f.Tags)
}

func TagsFromJSONB(data json.RawMessage) ([]string, error) {
	if len(data) == 0 {
		return []string{}, nil
	}
	var tags []string
	err := json.Unmarshal(data, &tags)
	return tags, err
}

