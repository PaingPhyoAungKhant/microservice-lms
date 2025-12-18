import type { User } from '../entities/User';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
  role?: 'student' | 'instructor';
}

export interface AuthTokens {
  accessToken: string;
  refreshToken?: string;
}

export interface AuthResult {
  user: User;
  tokens: AuthTokens;
}

export interface IAuthRepository {
  login(credentials: LoginCredentials): Promise<AuthResult>;
  register(data: RegisterData): Promise<AuthResult>;
  verifyToken(token: string): Promise<User>;
  refreshToken(refreshToken: string): Promise<AuthTokens>;
  logout(): Promise<void>;
}

