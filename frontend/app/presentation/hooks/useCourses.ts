import { useGetCoursesQuery, useGetCourseQuery } from '../../infrastructure/api/rtk/coursesApi';
import type { Course } from '../../domain/entities/Course';
import type { FetchBaseQueryError } from '@reduxjs/toolkit/query';

interface ApiException {
  message: string;
  status: number;
  code?: string;
  errors?: Record<string, string[]>;
}

export interface UseCoursesReturn {
  courses: Course[];
  total: number;
  loading: boolean;
  error: ApiException | null;
  fetchCourses: (params?: {
    searchQuery?: string;
    category?: string;
    instructorId?: string;
    limit?: number;
    offset?: number;
    sortColumn?: string;
    sortDirection?: 'asc' | 'desc';
  }) => void;
  refetch: () => void;
  reset: () => void;
}

export function useCourses(params?: {
  searchQuery?: string;
  category?: string;
  instructorId?: string;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}, options?: { skip?: boolean }): UseCoursesReturn {
  const apiParams = {
    searchQuery: params?.searchQuery,
    categoryId: params?.category,
    limit: params?.limit,
    offset: params?.offset,
    sortColumn: params?.sortColumn,
    sortDirection: params?.sortDirection,
  };

  const { data, isLoading, error, refetch } = useGetCoursesQuery(apiParams, {
    skip: options?.skip || false,
  });

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to fetch courses',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  return {
    courses: data?.courses || [],
    total: data?.total || 0,
    loading: isLoading,
    error: getError(),
    fetchCourses: () => {
    },
    refetch: () => {
      refetch();
    },
    reset: () => {
    },
  };
}

export interface UseCourseDetailReturn {
  course: Course | null;
  loading: boolean;
  error: ApiException | null;
  fetchCourse: (courseId: string) => void;
  refetch: () => void;
  reset: () => void;
}

export function useCourseDetail(courseId?: string): UseCourseDetailReturn {
  const { data, isLoading, error, refetch } = useGetCourseQuery(courseId || '', {
    skip: !courseId,
  });

  const getError = (): ApiException | null => {
    if (error) {
      const err = error as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Failed to fetch course',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  return {
    course: data || null,
    loading: isLoading,
    error: getError(),
    fetchCourse: () => {
    },
    refetch: () => {
      refetch();
    },
    reset: () => {
    },
  };
}
