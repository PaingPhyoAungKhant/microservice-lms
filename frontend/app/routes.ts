
import { type RouteConfig, route, index, prefix } from '@react-router/dev/routes';

export default [
  route('/', './routes/Index.tsx', [
    index('./routes/Home.tsx'),
    ...prefix('courses', [
      index('./routes/Courses.tsx'),
      route(':id', './routes/CourseDetail.tsx'),
    ]),
    ...prefix('auth', [
      // index('/', './routes/Login.tsx'),
      route('login', './routes/Login.tsx'),
      route('register', './routes/Register.tsx'),
    ]),
    route('register', './routes/RegisterRoot.tsx'),
    route('forgot-password', './routes/ForgotPassword.tsx'),
  ]),
  route('dashboard', './routes/dashboard/Index.tsx', [
    index('./routes/dashboard/Login.tsx'),
    ...prefix('admin', [
      route('users', './routes/dashboard/admin/Users.tsx'),
      route('users/:id', './routes/dashboard/admin/UserDetail.tsx'),
      route('categories', './routes/dashboard/admin/Categories.tsx'),
      route('courses', './routes/dashboard/admin/Courses.tsx'),
      route('courses/:id', './routes/dashboard/admin/CourseDetail.tsx'),
      route('course-offerings', './routes/dashboard/admin/CourseOfferings.tsx'),
      route('course-offerings/:id', './routes/dashboard/admin/CourseOfferingDetail.tsx'),
      route('enrollments', './routes/dashboard/admin/Enrollments.tsx'),
    ]),
    route('instructor', './routes/dashboard/Instructor.tsx', [
      index('./routes/dashboard/instructor/Index.tsx'),
      route('course-offerings', './routes/dashboard/instructor/CourseOfferings.tsx'),
      route('course-offerings/:id', './routes/dashboard/instructor/CourseOfferingDetail.tsx'),
      route('enrollments', './routes/dashboard/instructor/Enrollments.tsx'),
    ]),
    route('student', './routes/dashboard/Student.tsx', [
      index('./routes/dashboard/student/StudentDashboard.tsx'),
      route('classes/:offeringId', './routes/dashboard/student/ClassView.tsx'),
    ]),
  ]),
] satisfies RouteConfig;
