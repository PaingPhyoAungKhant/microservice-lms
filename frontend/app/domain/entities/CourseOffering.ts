export interface CourseOfferingInstructor {
  id: string;
  courseOfferingId: string;
  instructorId: string;
  instructorUsername: string;
  createdAt: string;
  updatedAt: string;
}

export interface CourseOffering {
  id: string;
  courseId: string;
  courseName?: string;
  name: string;
  description: string;
  offeringType: 'online' | 'oncampus';
  status: 'pending' | 'active' | 'ongoing' | 'completed';
  duration?: string;
  classTime?: string;
  enrollmentCost: number;
  createdAt: string;
  updatedAt: string;
  instructors?: CourseOfferingInstructor[];
}

