import Enrollments from '../shared/Enrollments';

export default function AdminEnrollments() {
  return <Enrollments routePrefix="/dashboard/admin" allowedRoles={['admin']} />;
}
