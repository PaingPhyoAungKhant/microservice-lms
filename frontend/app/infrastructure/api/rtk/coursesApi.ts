import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import { 
  courseResponseSchema, 
  coursesPaginatedResponseSchema,
  createCourseRequestSchema,
  updateCourseRequestSchema,
} from '../../validation/schemas';
import type { Course } from '../../../domain/entities/Course';
import type { Category } from '../../../domain/entities/Category';
import type { z } from 'zod';

function transformCourse(course: z.infer<typeof courseResponseSchema>): Course {
  return {
    id: course.id,
    name: course.name,
    description: course.description,
    thumbnailId: course.thumbnail_id ?? undefined,
    thumbnailUrl: course.thumbnail_url,
    categories: course.categories?.map((cat) => ({
      id: cat.id,
      name: cat.name,
      description: cat.description,
      createdAt: cat.created_at,
      updatedAt: cat.updated_at,
    })) as Category[] | undefined,
    createdAt: course.created_at,
    updatedAt: course.updated_at,
  };
}

interface GetCoursesParams {
  searchQuery?: string;
  categoryId?: string;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

interface CoursesListResponse {
  courses: Course[];
  total: number;
}

interface CreateCourseRequest {
  name: string;
  description?: string;
  thumbnailId?: string | null;
  categoryIds?: string[];
}

interface UpdateCourseRequest {
  name?: string;
  description?: string;
  thumbnailId?: string | null;
  categoryIds?: string[];
}

export const coursesApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getCourses: builder.query<CoursesListResponse, GetCoursesParams | void>({
            query: (params) => {
              const searchParams = new URLSearchParams();
        if (params?.searchQuery) searchParams.append('search', params.searchQuery);
        if (params?.categoryId) searchParams.append('category_id', params.categoryId);
              if (params?.limit) searchParams.append('limit', params.limit.toString());
              if (params?.offset) searchParams.append('offset', params.offset.toString());
              if (params?.sortColumn) searchParams.append('sort_column', params.sortColumn);
              if (params?.sortDirection) searchParams.append('sort_direction', params.sortDirection);

              const queryString = searchParams.toString();
        return {
          url: `${endpoints.courses.base}${queryString ? `?${queryString}` : ''}`,
          method: 'GET',
        };
            },
            transformResponse: (response: unknown) => {
              const validated = coursesPaginatedResponseSchema.parse(response);
              return {
          courses: validated.courses.map(transformCourse),
                total: validated.total,
              };
            },
      providesTags: (result) =>
        result
          ? [
              ...result.courses.map(({ id }) => ({ type: 'Course' as const, id })),
              { type: 'Course', id: 'LIST' },
            ]
          : [{ type: 'Course', id: 'LIST' }],
    }),

    getCourse: builder.query<Course, string>({
      query: (id) => ({
        url: endpoints.courses.byId(id),
        method: 'GET',
      }),
            transformResponse: (response: unknown) => {
              const validated = courseResponseSchema.parse(response);
        return transformCourse(validated);
            },
      providesTags: (_result, _error, id) => [{ type: 'Course', id }],
    }),

    createCourse: builder.mutation<Course, CreateCourseRequest>({
      query: (data) => ({
        url: endpoints.courses.base,
        method: 'POST',
        body: createCourseRequestSchema.parse({
          name: data.name,
          description: data.description,
          thumbnail_id: data.thumbnailId ?? null,
          category_ids: data.categoryIds,
        }),
      }),
      transformResponse: (response: unknown) => {
        const validated = courseResponseSchema.parse(response);
        return transformCourse(validated);
      },
      invalidatesTags: [{ type: 'Course', id: 'LIST' }],
    }),

    updateCourse: builder.mutation<Course, { id: string; data: UpdateCourseRequest }>({
      query: ({ id, data }) => ({
        url: endpoints.courses.byId(id),
        method: 'PUT',
        body: updateCourseRequestSchema.parse({
          name: data.name,
          description: data.description,
          thumbnail_id: data.thumbnailId ?? null,
          category_ids: data.categoryIds,
        }),
      }),
      transformResponse: (response: unknown) => {
        const validated = courseResponseSchema.parse(response);
        return transformCourse(validated);
      },
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'Course', id },
        { type: 'Course', id: 'LIST' },
      ],
    }),

    deleteCourse: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.courses.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'Course', id },
        { type: 'Course', id: 'LIST' },
      ],
    }),
  }),
});

export const {
  useGetCoursesQuery,
  useLazyGetCoursesQuery,
  useGetCourseQuery,
  useCreateCourseMutation,
  useUpdateCourseMutation,
  useDeleteCourseMutation,
} = coursesApi;

