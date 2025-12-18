import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import {
  courseSectionResponseSchema,
  courseSectionsResponseSchema,
  createCourseSectionRequestSchema,
  updateCourseSectionRequestSchema,
  reorderCourseSectionsRequestSchema,
} from '../../validation/schemas';
import type { CourseSection } from '../../../domain/entities/CourseSection';
import type { z } from 'zod';

function transformCourseSection(section: z.infer<typeof courseSectionResponseSchema>): CourseSection {
  return {
    id: section.id,
    courseOfferingId: section.course_offering_id,
    name: section.name,
    description: section.description,
    order: section.order,
    status: section.status,
    createdAt: section.created_at,
    updatedAt: section.updated_at,
  };
}

interface GetCourseSectionsParams {
  offeringId: string;
}

interface CourseSectionsListResponse {
  sections: CourseSection[];
}

interface CreateCourseSectionRequest {
  name: string;
  description?: string;
  order?: number;
  status?: 'draft' | 'published' | 'archived';
}

interface UpdateCourseSectionRequest {
  name?: string;
  description?: string;
  order?: number;
  status?: 'draft' | 'published' | 'archived';
}

interface ReorderItem {
  id: string;
  order: number;
}

export const courseSectionsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getCourseSections: builder.query<CourseSectionsListResponse, GetCourseSectionsParams>({
      query: ({ offeringId }) => ({
        url: endpoints.courseSections.byOffering(offeringId),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = courseSectionsResponseSchema.parse(response);
        return {
          sections: validated.sections.map(transformCourseSection),
        };
      },
      providesTags: (result, _error, { offeringId }) =>
        result
          ? [
              ...result.sections.map(({ id }) => ({ type: 'CourseSection' as const, id })),
              { type: 'CourseSection', id: `LIST-${offeringId}` },
            ]
          : [{ type: 'CourseSection', id: `LIST-${offeringId}` }],
    }),

    getCourseSection: builder.query<CourseSection, string>({
      query: (id) => ({
        url: endpoints.courseSections.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = courseSectionResponseSchema.parse(response);
        return transformCourseSection(validated);
      },
      providesTags: (_result, _error, id) => [{ type: 'CourseSection', id }],
    }),

    createCourseSection: builder.mutation<CourseSection, { offeringId: string; data: CreateCourseSectionRequest }>({
      query: ({ offeringId, data }) => ({
        url: endpoints.courseSections.byOffering(offeringId),
        method: 'POST',
        body: createCourseSectionRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = courseSectionResponseSchema.parse(response);
        return transformCourseSection(validated);
      },
      invalidatesTags: (_result, _error, { offeringId }) => [
        { type: 'CourseSection', id: 'LIST' },
        { type: 'CourseSection', id: `LIST-${offeringId}` },
      ],
    }),

    updateCourseSection: builder.mutation<CourseSection, { id: string; data: UpdateCourseSectionRequest }>({
      query: ({ id, data }) => ({
        url: endpoints.courseSections.byId(id),
        method: 'PUT',
        body: updateCourseSectionRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = courseSectionResponseSchema.parse(response);
        return transformCourseSection(validated);
      },
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'CourseSection', id },
        { type: 'CourseSection', id: 'LIST' },
      ],
    }),

    deleteCourseSection: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.courseSections.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'CourseSection', id },
        { type: 'CourseSection', id: 'LIST' },
      ],
    }),

    reorderCourseSections: builder.mutation<void, { offeringId: string; items: ReorderItem[] }>({
      query: ({ offeringId, items }) => ({
        url: endpoints.courseSections.reorder(offeringId),
        method: 'PUT',
        body: reorderCourseSectionsRequestSchema.parse({ items }),
      }),
      invalidatesTags: (_result, _error, { offeringId }) => [
        { type: 'CourseSection', id: 'LIST' },
        { type: 'CourseSection', id: `LIST-${offeringId}` },
      ],
    }),
  }),
});

export const {
  useGetCourseSectionsQuery,
  useLazyGetCourseSectionsQuery,
  useGetCourseSectionQuery,
  useCreateCourseSectionMutation,
  useUpdateCourseSectionMutation,
  useDeleteCourseSectionMutation,
  useReorderCourseSectionsMutation,
} = courseSectionsApi;

