import type { User, UserRole, UserStatus } from '../../domain/entities/User';

export interface UserDTO extends User {}

export interface UserListDTO {
  users: UserDTO[];
  total: number;
}

export interface CreateUserDTO {
  email: string;
  username: string;
  password: string;
  role: UserRole;
}

export interface UpdateUserDTO {
  username?: string;
  email?: string;
  role?: UserRole;
  status?: UserStatus;
}

