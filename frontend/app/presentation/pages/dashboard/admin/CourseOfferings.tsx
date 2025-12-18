import CourseOfferings from '../shared/CourseOfferings';

export default function AdminCourseOfferings() {
  return <CourseOfferings routePrefix="/dashboard/admin" allowedRoles={['admin']} />;
}
