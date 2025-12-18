export type EnrollmentStatus = 'pending' | 'approved' | 'rejected' | 'completed';

export interface Enrollment {
  id: string;
  studentId: string;
  studentUsername: string;
  courseId: string;
  courseName: string;
  courseOfferingId: string;
  courseOfferingName: string;
  status: EnrollmentStatus;
  createdAt: string;
  updatedAt: string;
}

