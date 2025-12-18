import CourseOfferingDetail from '../shared/CourseOfferingDetail';

export default function AdminCourseOfferingDetail() {
  return <CourseOfferingDetail routePrefix="/dashboard/admin" allowedRoles={['admin']} />;
}
