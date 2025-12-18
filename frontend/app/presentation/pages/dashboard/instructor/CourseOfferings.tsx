import CourseOfferings from '../shared/CourseOfferings';

export default function InstructorCourseOfferings() {
  return <CourseOfferings routePrefix="/dashboard/instructor" allowedRoles={['instructor', 'admin']} />;
}
