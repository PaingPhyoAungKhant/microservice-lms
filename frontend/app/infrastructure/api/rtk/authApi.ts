import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import {
  authResponseSchema,
  userResponseSchema,
  refreshTokenResponseSchema,
  registerRequestSchema,
  forgotPasswordRequestSchema,
  verifyOTPRequestSchema,
  verifyOTPResponseSchema,
  resetPasswordRequestSchema,
} from '../../validation/schemas';
import { storage } from '../../storage/storage';
import { setPublicUser, setDashboardUser, clearPublicAuth } from '../../store/authSlice';
import type { User } from '../../../domain/entities/User';
import type { z } from 'zod';


function transformUser(user: z.infer<typeof userResponseSchema>): User {
  return {
    id: user.id,
    username: user.username,
    email: user.email,
    role: user.role,
    status: user.status || 'pending',
    emailVerified: user.email_verified,
    createdAt: user.created_at || new Date().toISOString(),
    updatedAt: user.updated_at || new Date().toISOString(),
  };
}

interface LoginRequest {
  email: string;
  password: string;
}

interface AuthResponse {
  access_token: string;
  refresh_token?: string;
  user: User;
}

interface RefreshTokenRequest {
  refresh_token: string;
}

interface ForgotPasswordRequest {
  email: string;
}

interface VerifyOTPRequest {
  email: string;
  otp: string;
}

interface ResetPasswordRequest {
  token: string;
  new_password: string;
}


interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    login: builder.mutation<AuthResponse, LoginRequest>({
      query: (credentials) => ({
        url: endpoints.auth.login,
        method: 'POST',
        body: credentials,
      }),
      transformResponse: (response: unknown) => {
        const validated = authResponseSchema.parse(response);
        return {
          access_token: validated.access_token,
          refresh_token: validated.refresh_token,
          user: validated.user,
        };
      },
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          storage.setAccessToken(data.access_token);
          if (data.refresh_token) {
            storage.setRefreshToken(data.refresh_token);
          }
          storage.setUser(data.user);
          dispatch(setPublicUser(data.user));
        } catch (error) {
          const queryError = error as any;
          if (queryError?.status === 'FETCH_ERROR' || queryError?.status === 'PARSING_ERROR') {
            return;
          }
          console.error('Failed to store tokens after login:', error);
        }
      },
      invalidatesTags: ['Auth', 'User'],
    }),

    loginDashboard: builder.mutation<AuthResponse, LoginRequest>({
      query: (credentials) => ({
        url: endpoints.auth.login,
        method: 'POST',
        body: credentials,
      }),
      transformResponse: (response: unknown) => {
        const validated = authResponseSchema.parse(response);
        return {
          access_token: validated.access_token,
          refresh_token: validated.refresh_token,
          user: validated.user,
        };
      },
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          storage.setDashboardAccessToken(data.access_token);
          if (data.refresh_token) {
            storage.setDashboardRefreshToken(data.refresh_token);
          }
          storage.setDashboardUser(data.user);
          dispatch(setDashboardUser(data.user));
        } catch (error) {
          const queryError = error as any;
          if (queryError?.status === 'FETCH_ERROR' || queryError?.status === 'PARSING_ERROR') {
            return;
          }
          console.error('Failed to store tokens after dashboard login:', error);
        }
      },
      invalidatesTags: ['Auth', 'User'],
    }),

    register: builder.mutation<User, RegisterRequest>({
      query: (data) => ({
        url: endpoints.auth.register,
        method: 'POST',
        body: registerRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = userResponseSchema.parse(response);
        return transformUser(validated);
      },
      invalidatesTags: ['Auth', 'User'],
    }),

    verifyToken: builder.query<User, void>({
      queryFn: async (_arg, _queryApi, _extraOptions, baseQuery) => {
        try {
          const result = await baseQuery({
            url: endpoints.auth.verify,
            method: 'GET',
          });
          
          if (result.error) {
            if (result.error.status === 401 || result.error.status === 403) {
              storage.removeAccessToken();
              storage.removeRefreshToken();
              storage.removeUser();
            }
            return { error: result.error };
          }
          
          const meta = (result as any).meta;
          if (meta?.response) {
            const headers = meta.response.headers;
            if (headers instanceof Headers) {
              const userId = headers.get('X-User-ID');
              const userEmail = headers.get('X-User-Email');
              const userRole = headers.get('X-User-Role');
              
              if (userId && userEmail && userRole) {
                const user = transformUser({
                  id: userId,
                  email: userEmail,
                  username: userEmail.split('@')[0],
                  role: userRole as 'student' | 'instructor' | 'admin',
                  status: 'active' as const,
                  email_verified: true,
                  created_at: new Date().toISOString(),
                  updated_at: new Date().toISOString(),
                });
                return { data: user };
              }
            }
          }
          
          if (result.data && typeof result.data === 'object' && Object.keys(result.data as object).length > 0) {
            try {
              const validated = userResponseSchema.parse(result.data);
              return { data: transformUser(validated) };
            } catch (e) {
              console.error('Failed to transform user:', e);
            }
          }
          
          return {
            error: {
              status: 'CUSTOM_ERROR',
              error: 'Invalid verify response: no user data found in headers or body',
            },
          };
        } catch (error) {
          return {
            error: {
              status: 'FETCH_ERROR',
              error: 'Network error while verifying token',
            },
          };
        }
      },
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          if (data) {
            storage.setUser(data);
            dispatch(setPublicUser(data));
          }
        } catch (error) {
          const queryError = error as any;
          if (queryError?.status === 401 || queryError?.status === 403 || (queryError?.error as any)?.status === 401) {
            storage.removeAccessToken();
            storage.removeRefreshToken();
            storage.removeUser();
            dispatch(clearPublicAuth());
          }
        }
      },
      providesTags: ['Auth'],
    }),

    refreshToken: builder.mutation<{ access_token: string; refresh_token?: string }, RefreshTokenRequest>({
      query: (data) => ({
        url: endpoints.auth.refresh,
        method: 'POST',
        body: data,
      }),
      transformResponse: (response: unknown) => {
        const validated = refreshTokenResponseSchema.parse(response);
        return {
          access_token: validated.access_token,
          refresh_token: validated.refresh_token,
        };
      },
      async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled;
          const dashboardToken = storage.getDashboardRefreshToken();
          if (dashboardToken === _arg.refresh_token) {
            storage.setDashboardAccessToken(data.access_token);
            if (data.refresh_token) {
              storage.setDashboardRefreshToken(data.refresh_token);
            }
          } else {
            storage.setAccessToken(data.access_token);
            if (data.refresh_token) {
              storage.setRefreshToken(data.refresh_token);
            }
          }
        } catch (error) {
          console.error('Failed to store tokens after refresh:', error);
        }
      },
    }),

    forgotPassword: builder.mutation<{ message: string }, ForgotPasswordRequest>({
      query: (data) => ({
        url: endpoints.auth.forgotPassword,
        method: 'POST',
        body: forgotPasswordRequestSchema.parse(data),
      }),
    }),

    verifyOTP: builder.mutation<{ isValid: boolean; passwordResetToken?: string; errorMessage?: string }, VerifyOTPRequest>({
      query: (data) => ({
        url: endpoints.auth.verifyOTP,
        method: 'POST',
        body: verifyOTPRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = verifyOTPResponseSchema.parse(response);
        return {
          isValid: validated.IsValid,
          passwordResetToken: validated.PasswordResetToken,
          errorMessage: validated.ErrorMessage,
        };
      },
    }),

    resetPassword: builder.mutation<{ message: string }, ResetPasswordRequest>({
      query: (data) => ({
        url: endpoints.auth.resetPassword,
        method: 'POST',
        body: resetPasswordRequestSchema.parse(data),
      }),
    }),
  }),
});

export const {
  useLoginMutation,
  useLoginDashboardMutation,
  useRegisterMutation,
  useVerifyTokenQuery,
  useRefreshTokenMutation,
  useForgotPasswordMutation,
  useVerifyOTPMutation,
  useResetPasswordMutation,
} = authApi;

