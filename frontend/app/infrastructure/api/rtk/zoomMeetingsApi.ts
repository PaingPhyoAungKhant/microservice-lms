import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import {
  zoomMeetingResponseSchema,
  createZoomMeetingRequestSchema,
  updateZoomMeetingRequestSchema,
} from '../../validation/schemas';
import type { ZoomMeeting } from '../../../domain/entities/ZoomMeeting';
import type { z } from 'zod';

function transformZoomMeeting(meeting: z.infer<typeof zoomMeetingResponseSchema>): ZoomMeeting {
  return {
    id: meeting.id,
    sectionModuleId: meeting.section_module_id,
    zoomMeetingId: meeting.zoom_meeting_id,
    topic: meeting.topic,
    startTime: meeting.start_time ?? undefined,
    duration: meeting.duration ?? undefined,
    joinUrl: meeting.join_url,
    startUrl: meeting.start_url,
    password: meeting.password ?? undefined,
    createdAt: meeting.created_at,
    updatedAt: meeting.updated_at,
  };
}

interface CreateZoomMeetingRequest {
  sectionModuleId: string;
  topic: string;
  startTime?: string | null;
  duration?: number | null;
  password?: string | null;
}

interface UpdateZoomMeetingRequest {
  topic?: string;
  startTime?: string | null;
  duration?: number | null;
  password?: string | null;
}

function transformCreateMeetingRequest(data: CreateZoomMeetingRequest): z.infer<typeof createZoomMeetingRequestSchema> {
  return {
    section_module_id: data.sectionModuleId,
    topic: data.topic,
    start_time: data.startTime || null,
    duration: data.duration || null,
    password: data.password || null,
  };
}

function transformUpdateMeetingRequest(data: UpdateZoomMeetingRequest): z.infer<typeof updateZoomMeetingRequestSchema> {
  return {
    topic: data.topic,
    start_time: data.startTime || null,
    duration: data.duration || null,
    password: data.password || null,
  };
}

export const zoomMeetingsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getZoomMeeting: builder.query<ZoomMeeting, string>({
      query: (id) => ({
        url: endpoints.zoomMeetings.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = zoomMeetingResponseSchema.parse(response);
        return transformZoomMeeting(validated);
      },
      providesTags: (_result, _error, id) => [{ type: 'ZoomMeeting', id }],
    }),

    getZoomMeetingByModule: builder.query<ZoomMeeting, string>({
      query: (moduleId) => ({
        url: endpoints.zoomMeetings.byModule(moduleId),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = zoomMeetingResponseSchema.parse(response);
        return transformZoomMeeting(validated);
      },
      providesTags: (_result, _error, moduleId) => [{ type: 'ZoomMeeting', id: `MODULE-${moduleId}` }],
    }),

    createZoomMeeting: builder.mutation<ZoomMeeting, CreateZoomMeetingRequest>({
      query: (data) => {
        const transformedData = transformCreateMeetingRequest(data);
        return {
          url: endpoints.zoomMeetings.base,
          method: 'POST',
          body: createZoomMeetingRequestSchema.parse(transformedData),
        };
      },
      transformResponse: (response: unknown) => {
        const validated = zoomMeetingResponseSchema.parse(response);
        return transformZoomMeeting(validated);
      },
      invalidatesTags: (_result, _error, { sectionModuleId }) => [
        { type: 'ZoomMeeting', id: 'LIST' },
        { type: 'ZoomMeeting', id: `MODULE-${sectionModuleId}` },
        { type: 'SectionModule', id: sectionModuleId },
        { type: 'SectionModule', id: 'LIST' },
      ],
    }),

    updateZoomMeeting: builder.mutation<ZoomMeeting, { id: string; data: UpdateZoomMeetingRequest }>({
      query: ({ id, data }) => {
        const transformedData = transformUpdateMeetingRequest(data);
        return {
          url: endpoints.zoomMeetings.byId(id),
          method: 'PUT',
          body: updateZoomMeetingRequestSchema.parse(transformedData),
        };
      },
      transformResponse: (response: unknown) => {
        const validated = zoomMeetingResponseSchema.parse(response);
        return transformZoomMeeting(validated);
      },
      invalidatesTags: (result, _error, { id }) => {
        const tags = [
          { type: 'ZoomMeeting' as const, id },
          { type: 'ZoomMeeting' as const, id: 'LIST' },
        ];
        if (result?.sectionModuleId) {
          tags.push(
            { type: 'SectionModule' as const, id: result.sectionModuleId },
            { type: 'SectionModule' as const, id: 'LIST' }
          );
        }
        return tags;
      },
    }),

    deleteZoomMeeting: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.zoomMeetings.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'ZoomMeeting', id },
        { type: 'ZoomMeeting', id: 'LIST' },
        { type: 'SectionModule', id: 'LIST' },
      ],
    }),
  }),
});

export const {
  useGetZoomMeetingQuery,
  useLazyGetZoomMeetingQuery,
  useGetZoomMeetingByModuleQuery,
  useLazyGetZoomMeetingByModuleQuery,
  useCreateZoomMeetingMutation,
  useUpdateZoomMeetingMutation,
  useDeleteZoomMeetingMutation,
} = zoomMeetingsApi;

