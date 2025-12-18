package events

import "time"

type EnrollmentCreatedEvent struct {
	ID                 string    `json:"id"`
	StudentID          string    `json:"student_id"`
	StudentUsername    string    `json:"student_username"`
	CourseID           string    `json:"course_id"`
	CourseName         string    `json:"course_name"`
	CourseOfferingID   string    `json:"course_offering_id"`
	CourseOfferingName string    `json:"course_offering_name"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type EnrollmentUpdatedEvent struct {
	ID                 string    `json:"id"`
	StudentID          string    `json:"student_id"`
	StudentUsername    string    `json:"student_username"`
	CourseID           string    `json:"course_id"`
	CourseName         string    `json:"course_name"`
	CourseOfferingID   string    `json:"course_offering_id"`
	CourseOfferingName string    `json:"course_offering_name"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type EnrollmentDeletedEvent struct {
	ID                 string    `json:"id"`
	StudentID          string    `json:"student_id"`
	StudentUsername    string    `json:"student_username"`
	CourseID           string    `json:"course_id"`
	CourseName         string    `json:"course_name"`
	CourseOfferingID   string    `json:"course_offering_id"`
	CourseOfferingName string    `json:"course_offering_name"`
	DeletedAt          time.Time `json:"deleted_at"`
}

