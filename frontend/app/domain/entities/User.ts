export type UserRole = 'student' | 'instructor' | 'admin';
export type UserStatus = 'active' | 'inactive' | 'pending' | 'banned';

export interface User {
  id: string;
  username: string;
  email: string;
  role: UserRole;
  status?: UserStatus;
  emailVerified?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

