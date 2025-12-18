export const endpoints = {
  // Auth endpoints
  auth: {
    login: '/api/v1/auth/login',
    register: '/api/v1/auth/register',
    verify: '/api/v1/auth/verify',
    refresh: '/api/v1/auth/refresh-token',
    forgotPassword: '/api/v1/auth/forgot-password',
    verifyOTP: '/api/v1/auth/verify-otp',
    resetPassword: '/api/v1/auth/reset-password',
  },
  // User endpoints
  users: {
    base: '/api/v1/users',
    byId: (id: string) => `/api/v1/users/${id}`,
  },
  // Course endpoints
  courses: {
    base: '/api/v1/courses',
    byId: (id: string) => `/api/v1/courses/${id}`,
    find: '/api/v1/courses/find',
  },
  // Enrollment endpoints
  enrollments: {
    base: '/api/v1/enrollments',
    byId: (id: string) => `/api/v1/enrollments/${id}`,
    updateStatus: (id: string) => `/api/v1/enrollments/${id}/status`,
    byUser: (userId: string) => `/api/v1/enrollments/user/${userId}`,
    byCourse: (courseId: string) => `/api/v1/enrollments/course/${courseId}`,
  },
  // Category endpoints
  categories: {
    base: '/api/v1/categories',
    byId: (id: string) => `/api/v1/categories/${id}`,
  },
  // File endpoints
  files: {
    base: '/api/v1/files',
    byId: (id: string) => `/api/v1/files/${id}`,
    upload: '/api/v1/files',
    download: (id: string) => `/api/v1/files/${id}/download`,
  },
  // Bucket endpoints
  buckets: {
    download: (bucket: string, id: string) => `/api/v1/buckets/${bucket}/files/${id}/download`,
  },
  // Course Offering endpoints
  courseOfferings: {
    base: '/api/v1/course-offerings',
    byId: (id: string) => `/api/v1/course-offerings/${id}`,
    byCourse: (courseId: string) => `/api/v1/courses/${courseId}/offerings`,
    assignInstructor: (offeringId: string) => `/api/v1/course-offerings/${offeringId}/instructors`,
    removeInstructor: (offeringId: string, instructorId: string) => `/api/v1/course-offerings/${offeringId}/instructors/${instructorId}`,
  },
  // Course Section endpoints
  courseSections: {
    byOffering: (offeringId: string) => `/api/v1/course-offerings/${offeringId}/sections`,
    reorder: (offeringId: string) => `/api/v1/course-offerings/${offeringId}/sections/reorder`,
    byId: (id: string) => `/api/v1/course-sections/${id}`,
  },
  // Section Module endpoints
  sectionModules: {
    bySection: (sectionId: string) => `/api/v1/course-sections/${sectionId}/modules`,
    reorder: (sectionId: string) => `/api/v1/course-sections/${sectionId}/modules/reorder`,
    byId: (id: string) => `/api/v1/section-modules/${id}`,
  },
  // Zoom Meeting endpoints
  zoomMeetings: {
    base: '/api/v1/zoom/meetings',
    byId: (id: string) => `/api/v1/zoom/meetings/${id}`,
    byModule: (moduleId: string) => `/api/v1/zoom/meetings/module/${moduleId}`,
  },
  // Zoom Recording endpoints
  zoomRecordings: {
    base: '/api/v1/zoom/recordings',
    byId: (id: string) => `/api/v1/zoom/recordings/${id}`,
    byMeeting: (meetingId: string) => `/api/v1/zoom/recordings/meeting/${meetingId}`,
  },
} as const;

