import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import {
  courseOfferingResponseSchema,
  courseOfferingDetailResponseSchema,
  courseOfferingInstructorSchema,
  courseOfferingsPaginatedResponseSchema,
  createCourseOfferingRequestSchema,
  updateCourseOfferingRequestSchema,
} from '../../validation/schemas';
import type { CourseOffering, CourseOfferingInstructor } from '../../../domain/entities/CourseOffering';
import type { z } from 'zod';


function transformCourseOfferingInstructor(instructor: z.infer<typeof courseOfferingInstructorSchema>): CourseOfferingInstructor {
  return {
    id: instructor.id,
    courseOfferingId: instructor.course_offering_id,
    instructorId: instructor.instructor_id,
    instructorUsername: instructor.instructor_username,
    createdAt: instructor.created_at,
    updatedAt: instructor.updated_at,
  };
}

function transformCourseOffering(offering: z.infer<typeof courseOfferingResponseSchema>): CourseOffering {
  return {
    id: offering.id,
    courseId: offering.course_id,
    courseName: offering.course_name ?? undefined,
    name: offering.name,
    description: offering.description,
    offeringType: offering.offering_type,
    status: offering.status,
    duration: offering.duration ?? undefined,
    classTime: offering.class_time ?? undefined,
    enrollmentCost: offering.enrollment_cost,
    createdAt: offering.created_at,
    updatedAt: offering.updated_at,
  };
}

function transformCourseOfferingDetail(offering: z.infer<typeof courseOfferingDetailResponseSchema>): CourseOffering {
  return {
    ...transformCourseOffering(offering),
    instructors: offering.instructors?.map(transformCourseOfferingInstructor),
  };
}

interface GetCourseOfferingsParams {
  search?: string;
  courseId?: string;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

interface CourseOfferingsListResponse {
  courseOfferings: CourseOffering[];
  total: number;
}

interface CreateCourseOfferingRequest {
  name: string;
  description?: string;
  offeringType: 'online' | 'oncampus';
  duration?: string | null;
  classTime?: string | null;
  enrollmentCost: number;
}

interface UpdateCourseOfferingRequest {
  name?: string;
  description?: string;
  offeringType?: 'online' | 'oncampus';
  duration?: string | null;
  classTime?: string | null;
  enrollmentCost?: number;
  status?: 'pending' | 'active' | 'ongoing' | 'completed';
}

interface AssignInstructorRequest {
  instructor_id: string;
  instructor_username: string;
}

export const courseOfferingsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getCourseOfferings: builder.query<CourseOfferingsListResponse, GetCourseOfferingsParams | void>({
      query: (params) => {
        const searchParams = new URLSearchParams();
        if (params?.search) searchParams.append('search', params.search);
        if (params?.courseId) searchParams.append('course_id', params.courseId);
        if (params?.limit) searchParams.append('limit', params.limit.toString());
        if (params?.offset) searchParams.append('offset', params.offset.toString());
        if (params?.sortColumn) searchParams.append('sort_column', params.sortColumn);
        if (params?.sortDirection) searchParams.append('sort_direction', params.sortDirection);

        const queryString = searchParams.toString();
        return {
          url: `${endpoints.courseOfferings.base}${queryString ? `?${queryString}` : ''}`,
          method: 'GET',
        };
      },
      transformResponse: (response: unknown) => {
        const validated = courseOfferingsPaginatedResponseSchema.parse(response);
        return {
          courseOfferings: validated.offerings.map(transformCourseOffering),
          total: validated.total,
        };
      },
      providesTags: (result) =>
        result
          ? [
              ...result.courseOfferings.map(({ id }) => ({ type: 'CourseOffering' as const, id })),
              { type: 'CourseOffering', id: 'LIST' },
            ]
          : [{ type: 'CourseOffering', id: 'LIST' }],
    }),

    getCourseOffering: builder.query<CourseOffering, string>({
      query: (id) => ({
        url: endpoints.courseOfferings.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        try {
          const validated = courseOfferingDetailResponseSchema.parse(response);
          return transformCourseOfferingDetail(validated);
        } catch {
          const validated = courseOfferingResponseSchema.parse(response);
          return transformCourseOffering(validated);
        }
      },
      providesTags: (_result, _error, id) => [{ type: 'CourseOffering', id }],
    }),

    createCourseOffering: builder.mutation<CourseOffering, { courseId: string; data: CreateCourseOfferingRequest }>({
      query: ({ courseId, data }) => ({
        url: endpoints.courseOfferings.byCourse(courseId),
        method: 'POST',
        body: createCourseOfferingRequestSchema.parse({
          name: data.name,
          description: data.description,
          offering_type: data.offeringType,
          duration: data.duration ?? null,
          class_time: data.classTime ?? null,
          enrollment_cost: data.enrollmentCost,
        }),
      }),
      transformResponse: (response: unknown) => {
        const validated = courseOfferingResponseSchema.parse(response);
        return transformCourseOffering(validated);
      },
      invalidatesTags: [{ type: 'CourseOffering', id: 'LIST' }],
    }),

    updateCourseOffering: builder.mutation<CourseOffering, { id: string; data: UpdateCourseOfferingRequest }>({
      query: ({ id, data }) => ({
        url: endpoints.courseOfferings.byId(id),
        method: 'PUT',
        body: updateCourseOfferingRequestSchema.parse({
          name: data.name,
          description: data.description,
          offering_type: data.offeringType,
          duration: data.duration ?? null,
          class_time: data.classTime ?? null,
          enrollment_cost: data.enrollmentCost,
          status: data.status,
        }),
      }),
      transformResponse: (response: unknown) => {
        const validated = courseOfferingResponseSchema.parse(response);
        return transformCourseOffering(validated);
      },
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'CourseOffering', id },
        { type: 'CourseOffering', id: 'LIST' },
      ],
    }),

    assignInstructor: builder.mutation<void, { offeringId: string; data: AssignInstructorRequest }>({
      query: ({ offeringId, data }) => ({
        url: endpoints.courseOfferings.assignInstructor(offeringId),
        method: 'POST',
        body: {
          instructor_id: data.instructor_id,
          instructor_username: data.instructor_username,
        },
      }),
      invalidatesTags: (_result, _error, { offeringId }) => [{ type: 'CourseOffering', id: offeringId }],
    }),

    removeInstructor: builder.mutation<void, { offeringId: string; instructorId: string }>({
      query: ({ offeringId, instructorId }) => ({
        url: endpoints.courseOfferings.removeInstructor(offeringId, instructorId),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, { offeringId }) => [{ type: 'CourseOffering', id: offeringId }],
    }),

    deleteCourseOffering: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.courseOfferings.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'CourseOffering', id },
        { type: 'CourseOffering', id: 'LIST' },
      ],
    }),
  }),
});

export const {
  useGetCourseOfferingsQuery,
  useLazyGetCourseOfferingsQuery,
  useGetCourseOfferingQuery,
  useCreateCourseOfferingMutation,
  useUpdateCourseOfferingMutation,
  useDeleteCourseOfferingMutation,
  useAssignInstructorMutation,
  useRemoveInstructorMutation,
} = courseOfferingsApi;

