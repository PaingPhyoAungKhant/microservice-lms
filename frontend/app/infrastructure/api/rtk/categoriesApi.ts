import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import { categoryResponseSchema, categoriesResponseSchema, createCategoryRequestSchema, updateCategoryRequestSchema } from '../../validation/schemas';
import type { Category } from '../../../domain/entities/Category';
import type { z } from 'zod';


function transformCategory(category: z.infer<typeof categoryResponseSchema>): Category {
  return {
    id: category.id,
    name: category.name,
    description: category.description,
    createdAt: category.created_at,
    updatedAt: category.updated_at,
  };
}

interface GetCategoriesParams {
  search?: string;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

interface CategoriesListResponse {
  categories: Category[];
  total: number;
}

interface CreateCategoryRequest {
  name: string;
  description?: string;
}

interface UpdateCategoryRequest {
  name?: string;
  description?: string;
}

export const categoriesApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getCategories: builder.query<CategoriesListResponse, GetCategoriesParams | void>({
      query: (params) => {
        const searchParams = new URLSearchParams();
        if (params?.search) searchParams.append('search', params.search);
        if (params?.limit) searchParams.append('limit', params.limit.toString());
        if (params?.offset) searchParams.append('offset', params.offset.toString());
        if (params?.sortColumn) searchParams.append('sort_column', params.sortColumn);
        if (params?.sortDirection) searchParams.append('sort_direction', params.sortDirection);

        const queryString = searchParams.toString();
        return {
          url: `${endpoints.categories.base}${queryString ? `?${queryString}` : ''}`,
          method: 'GET',
        };
      },
      transformResponse: (response: unknown) => {
        const validated = categoriesResponseSchema.parse(response);
        return {
          categories: validated.categories.map(transformCategory),
          total: validated.total,
        };
      },
      providesTags: (result) =>
        result
          ? [
              ...result.categories.map(({ id }) => ({ type: 'Category' as const, id })),
              { type: 'Category', id: 'LIST' },
            ]
          : [{ type: 'Category', id: 'LIST' }],
    }),

    getCategory: builder.query<Category, string>({
      query: (id) => ({
        url: endpoints.categories.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = categoryResponseSchema.parse(response);
        return transformCategory(validated);
      },
      providesTags: (_result, _error, id) => [{ type: 'Category', id }],
    }),

    createCategory: builder.mutation<Category, CreateCategoryRequest>({
      query: (data) => ({
        url: endpoints.categories.base,
        method: 'POST',
        body: createCategoryRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = categoryResponseSchema.parse(response);
        return transformCategory(validated);
      },
      invalidatesTags: [{ type: 'Category', id: 'LIST' }],
    }),

    updateCategory: builder.mutation<Category, { id: string; data: UpdateCategoryRequest }>({
      query: ({ id, data }) => ({
        url: endpoints.categories.byId(id),
        method: 'PUT',
        body: updateCategoryRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = categoryResponseSchema.parse(response);
        return transformCategory(validated);
      },
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'Category', id },
        { type: 'Category', id: 'LIST' },
      ],
    }),

    deleteCategory: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.categories.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'Category', id },
        { type: 'Category', id: 'LIST' },
      ],
    }),
  }),
});

export const {
  useGetCategoriesQuery,
  useLazyGetCategoriesQuery,
  useGetCategoryQuery,
  useCreateCategoryMutation,
  useUpdateCategoryMutation,
  useDeleteCategoryMutation,
} = categoriesApi;

