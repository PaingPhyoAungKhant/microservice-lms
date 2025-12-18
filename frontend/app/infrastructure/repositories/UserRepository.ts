import type {
  IUserRepository,
  UserQuery,
  UserListResult,
  CreateUserData,
  UpdateUserData,
} from '../../domain/repositories/IUserRepository';
import type { User } from '../../domain/entities/User';
import { apiClient } from '../api/client';
import { endpoints } from '../api/endpoints';
import {
  userResponseSchema,
  usersResponseSchema,
  createUserRequestSchema,
  updateUserRequestSchema,
  type User as UserType,
} from '../validation/schemas';
import type { ApiException } from '../api/client';

interface UserListResponse {
  users: UserType[];
  total: number;
}

export class UserRepository implements IUserRepository {
  async findById(id: string): Promise<User | null> {
    try {
      const data = await apiClient.get<UserType>(endpoints.users.byId(id));
      const validated = userResponseSchema.parse(data);
      return this.mapToDomain(validated);
    } catch (error) {
      if (error instanceof ApiException && error.status === 404) {
        return null;
      }
      throw error;
    }
  }

  async find(query: UserQuery): Promise<UserListResult> {
    try {
      const params = new URLSearchParams();
      if (query.searchQuery) params.append('search_query', query.searchQuery);
      if (query.role) params.append('role', query.role);
      if (query.status) params.append('status', query.status);
      if (query.limit) params.append('limit', query.limit.toString());
      if (query.offset) params.append('offset', query.offset.toString());
      if (query.sortColumn) params.append('sort_column', query.sortColumn);
      if (query.sortDirection) params.append('sort_direction', query.sortDirection);

      const url = `${endpoints.users.base}?${params.toString()}`;
      const data = await apiClient.get<UserListResponse>(url);
      const validated = usersResponseSchema.parse(data);

      return {
        users: validated.users.map((user) => this.mapToDomain(user)),
        total: validated.total,
      };
    } catch (error) {
      throw error;
    }
  }

  async create(data: CreateUserData): Promise<User> {
    try {
      const validatedData = createUserRequestSchema.parse(data);
      const response = await apiClient.post<UserType>(endpoints.users.base, validatedData);
      const validated = userResponseSchema.parse(response);
      return this.mapToDomain(validated);
    } catch (error) {
      throw error;
    }
  }

  async update(id: string, data: UpdateUserData): Promise<User> {
    try {
      const validatedData = updateUserRequestSchema.parse(data);
      const response = await apiClient.put<UserType>(endpoints.users.byId(id), validatedData);
      const validated = userResponseSchema.parse(response);
      return this.mapToDomain(validated);
    } catch (error) {
      throw error;
    }
  }

  async delete(id: string): Promise<void> {
    try {
      await apiClient.delete<void>(endpoints.users.byId(id));
    } catch (error) {
      throw error;
    }
  }

  private mapToDomain(user: UserType): User {
    return {
      id: user.id,
      username: user.username,
      email: user.email,
      role: user.role,
      status: user.status,
      createdAt: user.created_at,
      updatedAt: user.updated_at,
    };
  }
}

