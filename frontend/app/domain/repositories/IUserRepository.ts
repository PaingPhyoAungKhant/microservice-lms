import type { User, UserRole, UserStatus } from '../entities/User';

export interface UserQuery {
  searchQuery?: string;
  role?: UserRole;
  status?: UserStatus;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

export interface UserListResult {
  users: User[];
  total: number;
}

export interface CreateUserData {
  email: string;
  username: string;
  password: string;
  role: UserRole;
}

export interface UpdateUserData {
  username?: string;
  email?: string;
  role?: UserRole;
  status?: UserStatus;
}

export interface IUserRepository {
  findById(id: string): Promise<User | null>;
  find(query: UserQuery): Promise<UserListResult>;
  create(data: CreateUserData): Promise<User>;
  update(id: string, data: UpdateUserData): Promise<User>;
  delete(id: string): Promise<void>;
}

