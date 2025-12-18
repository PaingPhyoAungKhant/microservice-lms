import CourseOfferingDetail from '../shared/CourseOfferingDetail';

export default function InstructorCourseOfferingDetail() {
  return <CourseOfferingDetail routePrefix="/dashboard/instructor" allowedRoles={['instructor', 'admin']} />;
}
