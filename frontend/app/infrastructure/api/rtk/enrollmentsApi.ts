import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import {
  enrollmentResponseSchema,
  enrollmentsPaginatedResponseSchema,
  createEnrollmentRequestSchema,
} from '../../validation/schemas';
import type { Enrollment, EnrollmentStatus } from '../../../domain/entities/Enrollment';
import type { z } from 'zod';

function transformEnrollment(enrollment: z.infer<typeof enrollmentResponseSchema>): Enrollment {
  return {
    id: enrollment.id,
    studentId: enrollment.student_id,
    studentUsername: enrollment.student_username,
    courseId: enrollment.course_id,
    courseName: enrollment.course_name,
    courseOfferingId: enrollment.course_offering_id,
    courseOfferingName: enrollment.course_offering_name,
    status: enrollment.status as EnrollmentStatus,
    createdAt: enrollment.created_at,
    updatedAt: enrollment.updated_at,
  };
}

interface GetEnrollmentsParams {
  searchQuery?: string;
  studentId?: string;
  courseId?: string;
  courseOfferingId?: string;
  status?: EnrollmentStatus;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

interface EnrollmentsListResponse {
  enrollments: Enrollment[];
  total: number;
}

interface CreateEnrollmentRequest {
  studentId: string;
  studentUsername: string;
  courseId: string;
  courseName: string;
  courseOfferingId: string;
  courseOfferingName: string;
}

export const enrollmentsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getEnrollments: builder.query<EnrollmentsListResponse, GetEnrollmentsParams | void>({
      query: (params) => {
        const searchParams = new URLSearchParams();
        if (params?.searchQuery) searchParams.append('search_query', params.searchQuery);
        if (params?.studentId) searchParams.append('student_id', params.studentId);
        if (params?.courseId) searchParams.append('course_id', params.courseId);
        if (params?.courseOfferingId) searchParams.append('course_offering_id', params.courseOfferingId);
        if (params?.status) searchParams.append('status', params.status);
        if (params?.limit !== undefined) searchParams.append('limit', params.limit.toString());
        if (params?.offset !== undefined) searchParams.append('offset', params.offset.toString());
        if (params?.sortColumn) searchParams.append('sort_column', params.sortColumn);
        if (params?.sortDirection) searchParams.append('sort_direction', params.sortDirection);

        const queryString = searchParams.toString();
        return {
          url: `${endpoints.enrollments.base}${queryString ? `?${queryString}` : ''}`,
          method: 'GET',
        };
      },
      transformResponse: (response: unknown) => {
        const validated = enrollmentsPaginatedResponseSchema.parse(response);
        return {
          enrollments: validated.enrollments.map(transformEnrollment),
          total: validated.total,
        };
      },
      providesTags: (result) =>
        result
          ? [
              ...result.enrollments.map(({ id }) => ({ type: 'Enrollment' as const, id })),
              { type: 'Enrollment', id: 'LIST' },
            ]
          : [{ type: 'Enrollment', id: 'LIST' }],
    }),

    getEnrollment: builder.query<Enrollment, string>({
      query: (id) => ({
        url: endpoints.enrollments.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = enrollmentResponseSchema.parse(response);
        return transformEnrollment(validated);
      },
      providesTags: (_result, _error, id) => [{ type: 'Enrollment', id }],
    }),

    createEnrollment: builder.mutation<Enrollment, CreateEnrollmentRequest>({
      query: (data) => ({
        url: endpoints.enrollments.base,
        method: 'POST',
        body: createEnrollmentRequestSchema.parse({
          student_id: data.studentId,
          student_username: data.studentUsername,
          course_id: data.courseId,
          course_name: data.courseName,
          course_offering_id: data.courseOfferingId,
          course_offering_name: data.courseOfferingName,
        }),
      }),
      transformResponse: (response: unknown) => {
        const validated = enrollmentResponseSchema.parse(response);
        return transformEnrollment(validated);
      },
      invalidatesTags: [{ type: 'Enrollment', id: 'LIST' }],
    }),

    updateEnrollmentStatus: builder.mutation<Enrollment, { id: string; status: EnrollmentStatus }>({
      query: ({ id, status }) => ({
        url: endpoints.enrollments.updateStatus(id),
        method: 'PUT',
        body: { status },
      }),
      transformResponse: (response: unknown) => {
        const validated = enrollmentResponseSchema.parse(response);
        return transformEnrollment(validated);
      },
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'Enrollment', id },
        { type: 'Enrollment', id: 'LIST' },
      ],
    }),

    deleteEnrollment: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.enrollments.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'Enrollment', id },
        { type: 'Enrollment', id: 'LIST' },
      ],
    }),
  }),
});

export const {
  useGetEnrollmentsQuery,
  useLazyGetEnrollmentsQuery,
  useGetEnrollmentQuery,
  useCreateEnrollmentMutation,
  useUpdateEnrollmentStatusMutation,
  useDeleteEnrollmentMutation,
} = enrollmentsApi;

