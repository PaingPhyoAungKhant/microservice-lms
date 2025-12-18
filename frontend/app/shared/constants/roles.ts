export type UserRole = 'student' | 'instructor' | 'admin';

export const ROLES = {
  STUDENT: 'student' as const,
  INSTRUCTOR: 'instructor' as const,
  ADMIN: 'admin' as const,
} as const;

export const ROLE_LABELS: Record<UserRole, string> = {
  student: 'Student',
  instructor: 'Instructor',
  admin: 'Admin',
};

