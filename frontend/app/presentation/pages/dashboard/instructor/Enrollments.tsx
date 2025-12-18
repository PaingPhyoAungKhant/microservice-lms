import Enrollments from '../shared/Enrollments';

export default function InstructorEnrollments() {
  return <Enrollments routePrefix="/dashboard/instructor" allowedRoles={['instructor', 'admin']} />;
}
