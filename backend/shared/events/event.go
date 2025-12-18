package events

const (
	// User Service events
	EventTypeUserCreated = "user.user.created"
	EventTypeUserUpdated = "user.user.updated"
	EventTypeUserDeleted = "user.user.deleted"

	// Auth Service events
	EventTypeAuthStudentRegistered  = "auth.student.registered"
	EventTypeAuthUserLoggedIn       = "auth.user.logged_in"
	EventTypeAuthUserLoggedOut      = "auth.user.logged_out"
	EventTypeAuthUserForgotPassword = "auth.user.forgot_password"
	EventTypeAuthUserResetPassword  = "auth.user.reset_password"
	EventTypeAuthUserRequestedEmailVerification = "auth.user.requested_email_verification"

	// Course Service events
	EventTypeCourseCreated                = "course.course.created"
	EventTypeCourseUpdated                = "course.course.updated"
	EventTypeCourseDeleted                = "course.course.deleted"
	EventTypeCourseOfferingCreated        = "course.offering.created"
	EventTypeCourseOfferingUpdated        = "course.offering.updated"
	EventTypeCourseOfferingDeleted        = "course.offering.deleted"
	EventTypeInstructorAssignedToOffering = "course.instructor.assigned"
	EventTypeInstructorRemovedFromOffering = "course.instructor.removed"
	EventTypeCourseSectionCreated         = "course.section.created"
	EventTypeCourseSectionUpdated         = "course.section.updated"
	EventTypeCourseSectionDeleted         = "course.section.deleted"
	EventTypeSectionModuleCreated         = "course.module.created"
	EventTypeSectionModuleUpdated         = "course.module.updated"
	EventTypeSectionModuleDeleted         = "course.module.deleted"

	// Zoom Service events
	EventTypeZoomMeetingCreated = "zoom.meeting.created"

	// Enrollment Service events
	EventTypeEnrollmentCreated = "enrollment.enrollment.created"
	EventTypeEnrollmentUpdated = "enrollment.enrollment.updated"
	EventTypeEnrollmentDeleted = "enrollment.enrollment.deleted"
)

