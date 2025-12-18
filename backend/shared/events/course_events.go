package events

import "time"



type CourseOfferingCreatedEvent struct {
	ID             string     `json:"id"`
	CourseID       string     `json:"course_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	OfferingType   string     `json:"offering_type"`
	Status         string     `json:"status"`
	Duration       *string    `json:"duration,omitempty"`
	ClassTime      *string    `json:"class_time,omitempty"`
	EnrollmentCost float64    `json:"enrollment_cost"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CourseOfferingUpdatedEvent struct {
	ID             string     `json:"id"`
	CourseID       string     `json:"course_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	OfferingType   string     `json:"offering_type"`
	Status         string     `json:"status"`
	Duration       *string    `json:"duration,omitempty"`
	ClassTime      *string    `json:"class_time,omitempty"`
	EnrollmentCost float64    `json:"enrollment_cost"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type InstructorAssignedToOfferingEvent struct {
	ID                 string    `json:"id"`
	CourseOfferingID   string    `json:"course_offering_id"`
	InstructorID       string    `json:"instructor_id"`
	InstructorUsername string    `json:"instructor_username"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type InstructorRemovedFromOfferingEvent struct {
	ID                 string    `json:"id"`
	CourseOfferingID   string    `json:"course_offering_id"`
	InstructorID       string    `json:"instructor_id"`
	InstructorUsername string    `json:"instructor_username"`
	RemovedAt          time.Time `json:"removed_at"`
}

type CourseSectionCreatedEvent struct {
	ID               string    `json:"id"`
	CourseOfferingID string    `json:"course_offering_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Order            int       `json:"order"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CourseSectionUpdatedEvent struct {
	ID               string    `json:"id"`
	CourseOfferingID string    `json:"course_offering_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Order            int       `json:"order"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type SectionModuleCreatedEvent struct {
	ID              string     `json:"id"`
	CourseSectionID string     `json:"course_section_id"`
	ContentID       *string    `json:"content_id,omitempty"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	ContentType     string     `json:"content_type"`
	ContentStatus   string     `json:"content_status"`
	Order           int        `json:"order"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type SectionModuleUpdatedEvent struct {
	ID              string     `json:"id"`
	CourseSectionID string     `json:"course_section_id"`
	ContentID       *string    `json:"content_id,omitempty"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	ContentType     string     `json:"content_type"`
	ContentStatus   string     `json:"content_status"`
	Order           int        `json:"order"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CourseCreatedEvent struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ThumbnailID *string    `json:"thumbnail_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CourseUpdatedEvent struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ThumbnailID *string    `json:"thumbnail_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CourseDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CourseOfferingDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CourseSectionDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

type SectionModuleDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

