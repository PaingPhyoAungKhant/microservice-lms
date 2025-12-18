import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import {
  zoomRecordingResponseSchema,
  zoomRecordingsResponseSchema,
  createZoomRecordingRequestSchema,
  updateZoomRecordingRequestSchema,
} from '../../validation/schemas';
import type { ZoomRecording } from '../../../domain/entities/ZoomRecording';
import type { z } from 'zod';

function transformZoomRecording(recording: z.infer<typeof zoomRecordingResponseSchema>): ZoomRecording {
  return {
    id: recording.id,
    zoomMeetingId: recording.zoom_meeting_id,
    fileId: recording.file_id,
    recordingType: recording.recording_type ?? undefined,
    recordingStartTime: recording.recording_start_time ?? undefined,
    recordingEndTime: recording.recording_end_time ?? undefined,
    fileSize: recording.file_size ?? undefined,
    createdAt: recording.created_at,
    updatedAt: recording.updated_at,
  };
}

interface CreateZoomRecordingRequest {
  zoomMeetingId: string;
  fileId: string;
  recordingType?: string | null;
  recordingStartTime?: string | null;
  recordingEndTime?: string | null;
  fileSize?: number | null;
}

interface UpdateZoomRecordingRequest {
  recordingType?: string | null;
  recordingStartTime?: string | null;
  recordingEndTime?: string | null;
  fileSize?: number | null;
}

function transformCreateRecordingRequest(data: CreateZoomRecordingRequest): z.infer<typeof createZoomRecordingRequestSchema> {
  if (!data.zoomMeetingId || typeof data.zoomMeetingId !== 'string') {
    throw new Error('zoomMeetingId is required and must be a string');
  }
  if (!data.fileId || typeof data.fileId !== 'string') {
    throw new Error('fileId is required and must be a string');
  }
  
  return {
    zoom_meeting_id: data.zoomMeetingId,
    file_id: data.fileId,
    recording_type: data.recordingType || null,
    recording_start_time: data.recordingStartTime || null,
    recording_end_time: data.recordingEndTime || null,
    file_size: data.fileSize || null,
  };
}

function transformUpdateRecordingRequest(data: UpdateZoomRecordingRequest): z.infer<typeof updateZoomRecordingRequestSchema> {
  return {
    recording_type: data.recordingType || null,
    recording_start_time: data.recordingStartTime || null,
    recording_end_time: data.recordingEndTime || null,
    file_size: data.fileSize || null,
  };
}

export const zoomRecordingsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getZoomRecording: builder.query<ZoomRecording, string>({
      query: (id) => ({
        url: endpoints.zoomRecordings.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = zoomRecordingResponseSchema.parse(response);
        return transformZoomRecording(validated);
      },
      providesTags: (_result, _error, id) => [{ type: 'ZoomRecording', id }],
    }),

    listZoomRecordings: builder.query<ZoomRecording[], string>({
      query: (meetingId) => ({
        url: endpoints.zoomRecordings.byMeeting(meetingId),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = zoomRecordingsResponseSchema.parse(response);
        return validated.map(transformZoomRecording);
      },
      providesTags: (result, _error, meetingId) =>
        result
          ? [
              ...result.map(({ id }) => ({ type: 'ZoomRecording' as const, id })),
              { type: 'ZoomRecording', id: `MEETING-${meetingId}` },
            ]
          : [{ type: 'ZoomRecording', id: `MEETING-${meetingId}` }],
    }),

    createZoomRecording: builder.mutation<ZoomRecording, CreateZoomRecordingRequest>({
      query: (data) => {
        console.log('[API] createZoomRecording - Input data:', {
          zoomMeetingId: data.zoomMeetingId,
          fileId: data.fileId,
          zoomMeetingIdType: typeof data.zoomMeetingId,
          fileIdType: typeof data.fileId,
        });
        
        const transformedData = transformCreateRecordingRequest(data);
        
        console.log('[API] createZoomRecording - Transformed data:', transformedData);
        
        return {
          url: endpoints.zoomRecordings.base,
          method: 'POST',
          body: createZoomRecordingRequestSchema.parse(transformedData),
        };
      },
      transformResponse: (response: unknown) => {
        const validated = zoomRecordingResponseSchema.parse(response);
        return transformZoomRecording(validated);
      },
      invalidatesTags: (_result, _error, { zoomMeetingId }) => [
        { type: 'ZoomRecording', id: 'LIST' },
        { type: 'ZoomRecording', id: `MEETING-${zoomMeetingId}` },
      ],
    }),

    updateZoomRecording: builder.mutation<ZoomRecording, { id: string; data: UpdateZoomRecordingRequest }>({
      query: ({ id, data }) => {
        const transformedData = transformUpdateRecordingRequest(data);
        return {
          url: endpoints.zoomRecordings.byId(id),
          method: 'PUT',
          body: updateZoomRecordingRequestSchema.parse(transformedData),
        };
      },
      transformResponse: (response: unknown) => {
        const validated = zoomRecordingResponseSchema.parse(response);
        return transformZoomRecording(validated);
      },
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'ZoomRecording', id },
        { type: 'ZoomRecording', id: 'LIST' },
      ],
    }),

    deleteZoomRecording: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.zoomRecordings.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'ZoomRecording', id },
        { type: 'ZoomRecording', id: 'LIST' },
      ],
    }),
  }),
});

export const {
  useGetZoomRecordingQuery,
  useLazyGetZoomRecordingQuery,
  useListZoomRecordingsQuery,
  useLazyListZoomRecordingsQuery,
  useCreateZoomRecordingMutation,
  useUpdateZoomRecordingMutation,
  useDeleteZoomRecordingMutation,
} = zoomRecordingsApi;

