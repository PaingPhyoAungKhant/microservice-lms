import { Navigate } from 'react-router';
import { ROUTES } from '../../../shared/constants/routes';

export default function InstructorIndex() {
  return <Navigate to={ROUTES.INSTRUCTOR_COURSE_OFFERINGS} replace />;
}

