import { useGetCategoriesQuery, useGetCategoryQuery, useCreateCategoryMutation, useUpdateCategoryMutation, useDeleteCategoryMutation } from '../../infrastructure/api/rtk/categoriesApi';
import type { Category } from '../../domain/entities/Category';
import type { FetchBaseQueryError } from '@reduxjs/toolkit/query';

interface ApiException {
  message: string;
  status: number;
  code?: string;
  errors?: Record<string, string[]>;
}

export interface UseCategoriesReturn {
  categories: Category[];
  total: number;
  loading: boolean;
  error: ApiException | null;
  fetchCategories: (params?: {
    search?: string;
    limit?: number;
    offset?: number;
    sortColumn?: string;
    sortDirection?: 'asc' | 'desc';
  }) => void;
  refetch: () => void;
  reset: () => void;
}

export function useCategories(params?: {
  search?: string;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}): UseCategoriesReturn {
  const { data, isLoading, error, refetch } = useGetCategoriesQuery(params || {}, {
    skip: false,
  });

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to fetch categories',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  return {
    categories: data?.categories || [],
    total: data?.total || 0,
    loading: isLoading,
    error: getError(),
    fetchCategories: () => {
    },
    refetch: () => {
      refetch();
    },
    reset: () => {
    },
  };
}

export interface UseCategoryDetailReturn {
  category: Category | null;
  loading: boolean;
  error: ApiException | null;
  fetchCategory: (categoryId: string) => void;
  refetch: () => void;
  reset: () => void;
}

export function useCategoryDetail(categoryId?: string): UseCategoryDetailReturn {
  const { data, isLoading, error, refetch } = useGetCategoryQuery(categoryId || '', {
    skip: !categoryId,
  });

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to fetch category',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  return {
    category: data || null,
    loading: isLoading,
    error: getError(),
    fetchCategory: () => {
    },
    refetch: () => {
      refetch();
    },
    reset: () => {
    },
  };
}

export interface UseCreateCategoryReturn {
  loading: boolean;
  error: ApiException | null;
  createCategory: (data: {
    name: string;
    description?: string;
  }) => Promise<Category>;
  reset: () => void;
}

export function useCreateCategory(): UseCreateCategoryReturn {
  const [createCategoryMutation, { isLoading, error }] = useCreateCategoryMutation();

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to create category',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  const createCategory = async (data: {
    name: string;
    description?: string;
  }): Promise<Category> => {
    try {
      const result = await createCategoryMutation(data).unwrap();
      return result;
    } catch (err) {
      throw err;
    }
  };

  return {
    loading: isLoading,
    error: getError(),
    createCategory,
    reset: () => {
    },
  };
}

export interface UseUpdateCategoryReturn {
  loading: boolean;
  error: ApiException | null;
  updateCategory: (
    categoryId: string,
    data: {
      name?: string;
      description?: string;
    }
  ) => Promise<Category>;
  reset: () => void;
}

export function useUpdateCategory(): UseUpdateCategoryReturn {
  const [updateCategoryMutation, { isLoading, error }] = useUpdateCategoryMutation();

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to update category',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  const updateCategory = async (
    categoryId: string,
    data: {
      name?: string;
      description?: string;
    }
  ): Promise<Category> => {
    try {
      const result = await updateCategoryMutation({ id: categoryId, data }).unwrap();
      return result;
    } catch (err) {
      throw err;
    }
  };

  return {
    loading: isLoading,
    error: getError(),
    updateCategory,
    reset: () => {
    },
  };
}

export interface UseDeleteCategoryReturn {
  loading: boolean;
  error: ApiException | null;
  deleteCategory: (categoryId: string) => Promise<void>;
  reset: () => void;
}

export function useDeleteCategory(): UseDeleteCategoryReturn {
  const [deleteCategoryMutation, { isLoading, error }] = useDeleteCategoryMutation();

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to delete category',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  const deleteCategory = async (categoryId: string): Promise<void> => {
    try {
      await deleteCategoryMutation(categoryId).unwrap();
    } catch (err) {
      throw err;
    }
  };

  return {
    loading: isLoading,
    error: getError(),
    deleteCategory,
    reset: () => {
    },
  };
}

