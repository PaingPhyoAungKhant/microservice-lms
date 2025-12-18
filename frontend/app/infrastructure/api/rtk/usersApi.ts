import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import { userResponseSchema, usersResponseSchema, createUserRequestSchema, updateUserRequestSchema } from '../../validation/schemas';
import type { User, UserRole, UserStatus } from '../../../domain/entities/User';
import type { z } from 'zod';


function transformUser(user: z.infer<typeof userResponseSchema>): User {
  return {
    id: user.id,
    username: user.username,
    email: user.email,
    role: user.role,
    status: user.status,
    emailVerified: user.email_verified,
    createdAt: user.created_at,
    updatedAt: user.updated_at,
  };
}

interface GetUsersParams {
  searchQuery?: string;
  role?: UserRole;
  status?: UserStatus;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

interface UsersListResponse {
  users: User[];
  total: number;
}

interface CreateUserRequest {
  email: string;
  username: string;
  password: string;
  role: UserRole;
}

interface UpdateUserRequest {
  username?: string;
  email?: string;
  role?: UserRole;
  status?: UserStatus;
}

export const usersApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getUsers: builder.query<UsersListResponse, GetUsersParams | void>({
      query: (params) => {
        const searchParams = new URLSearchParams();
        if (params?.searchQuery) searchParams.append('search_query', params.searchQuery);
        if (params?.role) searchParams.append('role', params.role);
        if (params?.status) searchParams.append('status', params.status);
        if (params?.limit !== undefined) searchParams.append('limit', params.limit.toString());
        if (params?.offset !== undefined) searchParams.append('offset', params.offset.toString());
        if (params?.sortColumn) searchParams.append('sort_column', params.sortColumn);
        if (params?.sortDirection) searchParams.append('sort_direction', params.sortDirection);

        const queryString = searchParams.toString();
        return {
          url: `${endpoints.users.base}${queryString ? `?${queryString}` : ''}`,
          method: 'GET',
        };
      },
      transformResponse: (response: unknown) => {
        const validated = usersResponseSchema.parse(response);
        return {
          users: validated.users.map(transformUser),
          total: validated.total,
        };
      },
      providesTags: (result) =>
        result
          ? [
              ...result.users.map(({ id }) => ({ type: 'User' as const, id })),
              { type: 'User', id: 'LIST' },
            ]
          : [{ type: 'User', id: 'LIST' }],
    }),

    getUser: builder.query<User, string>({
      query: (id) => ({
        url: endpoints.users.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = userResponseSchema.parse(response);
        return transformUser(validated);
      },
      providesTags: (_result, _error, id) => [{ type: 'User', id }],
    }),

    createUser: builder.mutation<User, CreateUserRequest>({
      query: (data) => ({
        url: endpoints.users.base,
        method: 'POST',
        body: createUserRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = userResponseSchema.parse(response);
        return transformUser(validated);
      },
      invalidatesTags: [{ type: 'User', id: 'LIST' }],
    }),

    updateUser: builder.mutation<User, { id: string; data: UpdateUserRequest }>({
      query: ({ id, data }) => ({
        url: endpoints.users.byId(id),
        method: 'PUT',
        body: updateUserRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = userResponseSchema.parse(response);
        return transformUser(validated);
      },
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'User', id },
        { type: 'User', id: 'LIST' },
      ],
    }),

    deleteUser: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.users.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'User', id },
        { type: 'User', id: 'LIST' },
      ],
    }),
  }),
});

export const {
  useGetUsersQuery,
  useLazyGetUsersQuery,
  useGetUserQuery,
  useCreateUserMutation,
  useUpdateUserMutation,
  useDeleteUserMutation,
} = usersApi;

