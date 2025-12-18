import { z } from 'zod';

// User schemas
export const userSchema = z.object({
  id: z.string(),
  username: z.string(),
  email: z.string().email(),
  role: z.enum(['student', 'instructor', 'admin']),
  status: z.enum(['active', 'inactive', 'pending', 'banned']).optional(),
  email_verified: z.boolean().optional(),
  email_verified_at: z.string().nullable().optional(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export const userResponseSchema = userSchema;
export const usersResponseSchema = z.object({
  users: z.array(userSchema),
  total: z.number(),
});

export const createUserRequestSchema = z.object({
  email: z.string().email(),
  username: z.string().min(3).max(255),
  password: z.string().min(8).max(255),
  role: z.enum(['student', 'instructor', 'admin']),
});

export const updateUserRequestSchema = z.object({
  username: z.string().min(3).max(255).optional(),
  email: z.string().email().optional(),
  role: z.enum(['student', 'instructor', 'admin']).optional(),
  status: z.enum(['active', 'inactive', 'pending', 'banned']).optional(),
});

// Course schemas
export const courseSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  thumbnail_id: z.string().nullable().optional(),
  thumbnail_url: z.string().optional(),
  categories: z.array(z.object({
    id: z.string(),
    name: z.string(),
    description: z.string(),
    created_at: z.string(),
    updated_at: z.string(),
  })).optional(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const courseResponseSchema = courseSchema;
export const coursesResponseSchema = z.array(courseSchema);
export const coursesPaginatedResponseSchema = z.object({
  courses: z.array(courseSchema),
  total: z.number(),
});

export const createCourseRequestSchema = z.object({
  name: z.string().min(1),
  description: z.string().optional(),
  thumbnail_id: z.string().nullable().optional(),
  category_ids: z.array(z.string()).optional(),
});

export const updateCourseRequestSchema = z.object({
  name: z.string().min(1).optional(),
  description: z.string().optional(),
  thumbnail_id: z.string().nullable().optional(),
  category_ids: z.array(z.string()).optional(),
});

// Auth schemas
export const loginRequestSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
});

export const registerRequestSchema = z.object({
  username: z.string().min(3).max(255),
  email: z.string().email(),
  password: z.string().min(8).max(255),
});

export const forgotPasswordRequestSchema = z.object({
  email: z.string().email(),
});

export const verifyOTPRequestSchema = z.object({
  email: z.string().email(),
  otp: z.string().min(6).max(6),
});

export const verifyOTPResponseSchema = z.object({
  IsValid: z.boolean(),
  PasswordResetToken: z.string().optional(),
  ErrorMessage: z.string().optional(),
}).passthrough();

export const resetPasswordRequestSchema = z.object({
  token: z.string().min(1),
  new_password: z.string().min(8).max(255),
});

export const authResponseSchema = z.object({
  access_token: z.string(),
  refresh_token: z.string().optional(),
  user: userSchema,
});

export const refreshTokenResponseSchema = z.object({
  access_token: z.string(),
  refresh_token: z.string().optional(),
});

// Enrollment schemas
export const enrollmentSchema = z.object({
  id: z.string(),
  student_id: z.string(),
  student_username: z.string(),
  course_id: z.string(),
  course_name: z.string(),
  course_offering_id: z.string(),
  course_offering_name: z.string(),
  status: z.enum(['pending', 'approved', 'rejected', 'completed']),
  created_at: z.string(),
  updated_at: z.string(),
});

export const enrollmentResponseSchema = enrollmentSchema;
export const enrollmentsResponseSchema = z.array(enrollmentSchema);
export const enrollmentsPaginatedResponseSchema = z.object({
  enrollments: z.array(enrollmentSchema),
  total: z.number(),
});

export const createEnrollmentRequestSchema = z.object({
  student_id: z.string().uuid(),
  student_username: z.string().min(3).max(255),
  course_id: z.string().uuid(),
  course_name: z.string().min(1).max(255),
  course_offering_id: z.string().uuid(),
  course_offering_name: z.string().min(1).max(255),
});

// Category schemas
export const categorySchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const categoryResponseSchema = categorySchema;
export const categoriesResponseSchema = z.object({
  categories: z.array(categorySchema),
  total: z.number(),
});

export const createCategoryRequestSchema = z.object({
  name: z.string().min(1),
  description: z.string().optional(),
});

export const updateCategoryRequestSchema = z.object({
  name: z.string().min(1).optional(),
  description: z.string().optional(),
});

// Course Offering schemas
export const courseOfferingSchema = z.object({
  id: z.string(),
  course_id: z.string(),
  course_name: z.string().nullable().optional(),
  name: z.string(),
  description: z.string(),
  offering_type: z.enum(['online', 'oncampus']),
  status: z.enum(['pending', 'active', 'ongoing', 'completed']),
  duration: z.string().nullable().optional(),
  class_time: z.string().nullable().optional(),
  enrollment_cost: z.number(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const courseOfferingInstructorSchema = z.object({
  id: z.string(),
  course_offering_id: z.string(),
  instructor_id: z.string(),
  instructor_username: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type CourseOfferingInstructorInput = z.infer<typeof courseOfferingInstructorSchema>;

export const courseOfferingResponseSchema = courseOfferingSchema;
export const courseOfferingDetailResponseSchema = courseOfferingSchema.extend({
  instructors: z.array(courseOfferingInstructorSchema),
  sections: z.array(z.any()),
});
export const courseOfferingsPaginatedResponseSchema = z.object({
  offerings: z.array(courseOfferingSchema),
  total: z.number(),
});

export const createCourseOfferingRequestSchema = z.object({
  name: z.string().min(1),
  description: z.string().optional(),
  offering_type: z.enum(['online', 'oncampus']),
  duration: z.string().nullable().optional(),
  class_time: z.string().nullable().optional(),
  enrollment_cost: z.number().min(0),
});

export const updateCourseOfferingRequestSchema = z.object({
  name: z.string().min(1).optional(),
  description: z.string().optional(),
  offering_type: z.enum(['online', 'oncampus']).optional(),
  duration: z.string().nullable().optional(),
  class_time: z.string().nullable().optional(),
  enrollment_cost: z.number().min(0).optional(),
  status: z.enum(['pending', 'active', 'ongoing', 'completed']).optional(),
});

// Course Section schemas
export const courseSectionSchema = z.object({
  id: z.string(),
  course_offering_id: z.string(),
  name: z.string(),
  description: z.string(),
  order: z.number(),
  status: z.enum(['draft', 'published', 'archived']),
  created_at: z.string(),
  updated_at: z.string(),
});

export const courseSectionResponseSchema = courseSectionSchema;
export const courseSectionsResponseSchema = z.object({
  sections: z.array(courseSectionSchema),
});

export const createCourseSectionRequestSchema = z.object({
  name: z.string().min(1),
  description: z.string().optional(),
  order: z.number().optional(),
  status: z.enum(['draft', 'published', 'archived']).optional(),
});

export const updateCourseSectionRequestSchema = z.object({
  name: z.string().min(1).optional(),
  description: z.string().optional(),
  order: z.number().optional(),
  status: z.enum(['draft', 'published', 'archived']).optional(),
});

export const reorderCourseSectionsRequestSchema = z.object({
  items: z.array(z.object({
    id: z.string(),
    order: z.number(),
  })),
});

// Section Module schemas
export const sectionModuleSchema = z.object({
  id: z.string(),
  course_section_id: z.string(),
  content_id: z.string().nullable().optional(),
  name: z.string(),
  description: z.string(),
  content_type: z.enum(['zoom']),
  content_status: z.enum(['draft', 'pending', 'created']),
  order: z.number(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const sectionModuleResponseSchema = sectionModuleSchema;
export const sectionModulesResponseSchema = z.object({
  modules: z.array(sectionModuleSchema),
});

export const createSectionModuleRequestSchema = z.object({
  name: z.string().min(1),
  description: z.string().optional(),
  content_type: z.enum(['zoom']),
  order: z.number().optional(),
});

export const updateSectionModuleRequestSchema = z.object({
  name: z.string().min(1).optional(),
  description: z.string().optional(),
  order: z.number().optional(),
});

export const reorderSectionModulesRequestSchema = z.object({
  items: z.array(z.object({
    id: z.string(),
    order: z.number(),
  })),
});

// Zoom Meeting schemas
export const zoomMeetingSchema = z.object({
  id: z.string(),
  section_module_id: z.string(),
  zoom_meeting_id: z.string(),
  topic: z.string(),
  start_time: z.string().nullable().optional(),
  duration: z.number().nullable().optional(),
  join_url: z.string(),
  start_url: z.string(),
  password: z.string().nullable().optional(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const zoomMeetingResponseSchema = zoomMeetingSchema;

export const createZoomMeetingRequestSchema = z.object({
  section_module_id: z.string(),
  topic: z.string().min(1),
  start_time: z.string().nullable().optional(),
  duration: z.number().nullable().optional(),
  password: z.string().nullable().optional(),
});

export const updateZoomMeetingRequestSchema = z.object({
  topic: z.string().min(1).optional(),
  start_time: z.string().nullable().optional(),
  duration: z.number().nullable().optional(),
  password: z.string().nullable().optional(),
});

// Zoom Recording schemas
export const zoomRecordingSchema = z.object({
  id: z.string(),
  zoom_meeting_id: z.string(),
  file_id: z.string(),
  recording_type: z.string().nullable().optional(),
  recording_start_time: z.string().nullable().optional(),
  recording_end_time: z.string().nullable().optional(),
  file_size: z.number().nullable().optional(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const zoomRecordingResponseSchema = zoomRecordingSchema;
export const zoomRecordingsResponseSchema = z.array(zoomRecordingSchema);

export const createZoomRecordingRequestSchema = z.object({
  zoom_meeting_id: z.string(),
  file_id: z.string(),
  recording_type: z.string().nullable().optional(),
  recording_start_time: z.string().nullable().optional(),
  recording_end_time: z.string().nullable().optional(),
  file_size: z.number().nullable().optional(),
});

export const updateZoomRecordingRequestSchema = z.object({
  recording_type: z.string().nullable().optional(),
  recording_start_time: z.string().nullable().optional(),
  recording_end_time: z.string().nullable().optional(),
  file_size: z.number().nullable().optional(),
});

// API Error schema
export const apiErrorSchema = z.object({
  message: z.string(),
  code: z.string().optional(),
  errors: z.record(z.string(), z.array(z.string())).optional(),
});

export type User = z.infer<typeof userSchema>;
export type Course = z.infer<typeof courseSchema>;
export type LoginRequest = z.infer<typeof loginRequestSchema>;
export type RegisterRequest = z.infer<typeof registerRequestSchema>;
export type AuthResponse = z.infer<typeof authResponseSchema>;
export type Enrollment = z.infer<typeof enrollmentSchema>;
export type ApiError = z.infer<typeof apiErrorSchema>;

