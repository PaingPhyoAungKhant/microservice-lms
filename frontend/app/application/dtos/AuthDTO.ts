import type { User } from '../../domain/entities/User';

export interface LoginDTO {
  email: string;
  password: string;
}

export interface RegisterDTO {
  name: string;
  email: string;
  password: string;
  role?: 'student' | 'instructor';
}

export interface AuthResponseDTO {
  user: User;
  accessToken: string;
  refreshToken?: string;
}

