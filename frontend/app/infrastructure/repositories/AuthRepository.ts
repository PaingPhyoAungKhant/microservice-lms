import type {
  IAuthRepository,
  LoginCredentials,
  RegisterData,
  AuthResult,
  AuthTokens,
} from '../../domain/repositories/IAuthRepository';
import type { User } from '../../domain/entities/User';
import { apiClient } from '../api/client';
import { endpoints } from '../api/endpoints';
import {
  authResponseSchema,
  userSchema,
  type AuthResponse as AuthResponseType,
  type User as UserType,
} from '../validation/schemas';

export class AuthRepository implements IAuthRepository {
  async login(credentials: LoginCredentials): Promise<AuthResult> {
    try {
      const data = await apiClient.post<AuthResponseType>(endpoints.auth.login, credentials);
      const validated = authResponseSchema.parse(data);

      return {
        user: this.mapUserToDomain(validated.user),
        tokens: {
          accessToken: validated.access_token ?? '',
          refreshToken: validated.refresh_token ?? '',
        },
      };
    } catch (error) {
      throw error;
    }
  }

  async register(data: RegisterData): Promise<AuthResult> {
    try {
      const response = await apiClient.post<AuthResponseType>(endpoints.auth.register, data);
      const validated = authResponseSchema.parse(response);

      return {
        user: this.mapUserToDomain(validated.user),
        tokens: {
          accessToken: validated.access_token ?? '',
          refreshToken: validated.refresh_token ?? '',
        },
      };
    } catch (error) {
      throw error;
    }
  }

  async verifyToken(token: string): Promise<User> {
    try {
      const data = await apiClient.post<UserType>(endpoints.auth.verify, { token });
      const validated = userSchema.parse(data);
      return this.mapUserToDomain(validated);
    } catch (error) {
      throw error;
    }
  }

  async refreshToken(refreshToken: string): Promise<AuthTokens> {
    try {
      const data = await apiClient.post<AuthResponseType>(endpoints.auth.refresh, {
        refreshToken,
      });
      const validated = authResponseSchema.parse(data);

      return {
        accessToken: validated.access_token ?? '',
        refreshToken: validated.refresh_token ?? '',
      };
    } catch (error) {
      throw error;
    }
  }

  async logout(): Promise<void> {
    return Promise.resolve();
  }

  private mapUserToDomain(user: UserType): User {
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

