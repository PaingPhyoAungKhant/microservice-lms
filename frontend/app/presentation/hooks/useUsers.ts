import { useGetUsersQuery, useGetUserQuery, useCreateUserMutation, useUpdateUserMutation, useDeleteUserMutation } from '../../infrastructure/api/rtk/usersApi';
import type { User, UserRole, UserStatus } from '../../domain/entities/User';
import type { FetchBaseQueryError } from '@reduxjs/toolkit/query';

interface ApiException {
  message: string;
  status: number;
  code?: string;
  errors?: Record<string, string[]>;
}

export interface UseUsersReturn {
  users: User[];
  total: number;
  loading: boolean;
  error: ApiException | null;
  fetchUsers: (params?: {
    searchQuery?: string;
    role?: UserRole;
    status?: UserStatus;
    limit?: number;
    offset?: number;
    sortColumn?: string;
    sortDirection?: 'asc' | 'desc';
  }) => void;
  refetch: () => void;
  reset: () => void;
}

export function useUsers(params?: {
  searchQuery?: string;
  role?: UserRole;
  status?: UserStatus;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}): UseUsersReturn {
  const { data, isLoading, error, refetch } = useGetUsersQuery(params || {}, {
    skip: false,
  });

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to fetch users',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  return {
    users: data?.users || [],
    total: data?.total || 0,
    loading: isLoading,
    error: getError(),
    fetchUsers: () => {
    },
    refetch: () => {
      refetch();
    },
    reset: () => {
    },
  };
}

export interface UseUserDetailReturn {
  user: User | null;
  loading: boolean;
  error: ApiException | null;
  fetchUser: (userId: string) => void;
  refetch: () => void;
  reset: () => void;
}

export function useUserDetail(userId?: string): UseUserDetailReturn {
  const { data, isLoading, error, refetch } = useGetUserQuery(userId || '', {
    skip: !userId,
  });

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to fetch user',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  return {
    user: data || null,
    loading: isLoading,
    error: getError(),
    fetchUser: () => {
    },
    refetch: () => {
      refetch();
    },
    reset: () => {
    },
  };
}

export interface UseCreateUserReturn {
  loading: boolean;
  error: ApiException | null;
  createUser: (data: {
    email: string;
    username: string;
    password: string;
    role: UserRole;
  }) => Promise<User>;
  reset: () => void;
}

export function useCreateUser(): UseCreateUserReturn {
  const [createUserMutation, { isLoading, error }] = useCreateUserMutation();

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to create user',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  const createUser = async (data: {
    email: string;
    username: string;
    password: string;
    role: UserRole;
  }): Promise<User> => {
    try {
      const result = await createUserMutation(data).unwrap();
      return result;
    } catch (err) {
      throw err;
    }
  };

  return {
    loading: isLoading,
    error: getError(),
    createUser,
    reset: () => {
    },
  };
}

export interface UseUpdateUserReturn {
  loading: boolean;
  error: ApiException | null;
  updateUser: (
    userId: string,
    data: {
      username?: string;
      email?: string;
      role?: UserRole;
      status?: UserStatus;
    }
  ) => Promise<User>;
  reset: () => void;
}

export function useUpdateUser(): UseUpdateUserReturn {
  const [updateUserMutation, { isLoading, error }] = useUpdateUserMutation();

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to update user',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  const updateUser = async (
    userId: string,
    data: {
      username?: string;
      email?: string;
      role?: UserRole;
      status?: UserStatus;
    }
  ): Promise<User> => {
    try {
      const result = await updateUserMutation({ id: userId, data }).unwrap();
      return result;
    } catch (err) {
      throw err;
    }
  };

  return {
    loading: isLoading,
    error: getError(),
    updateUser,
    reset: () => {
    },
  };
}

export interface UseDeleteUserReturn {
  loading: boolean;
  error: ApiException | null;
  deleteUser: (userId: string) => Promise<void>;
  reset: () => void;
}

export function useDeleteUser(): UseDeleteUserReturn {
  const [deleteUserMutation, { isLoading, error }] = useDeleteUserMutation();

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to delete user',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  const deleteUser = async (userId: string): Promise<void> => {
    try {
      await deleteUserMutation(userId).unwrap();
    } catch (err) {
      throw err;
    }
  };

  return {
    loading: isLoading,
    error: getError(),
    deleteUser,
    reset: () => {
    },
  };
}
