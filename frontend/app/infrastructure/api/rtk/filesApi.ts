import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import type { File } from '../../../domain/entities/File';
import { storage } from '../../storage/storage';

interface ListFilesParams {
  uploadedBy?: string;
  tags?: string[];
  mimeType?: string;
  bucketName?: string;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

interface FilesListResponse {
  files: File[];
  total: number;
}

interface UploadFileParams {
  file: File | Blob;
  bucketName?: string;
  tags?: string[];
}

interface FileResponse {
  id: string;
  original_filename: string;
  stored_filename: string;
  bucket_name: string;
  mime_type: string;
  size_bytes: number;
  uploaded_by: string;
  tags: string[];
  download_url?: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

function transformFile(file: FileResponse, apiGatewayURL?: string): File {
  return {
    id: file.id,
    originalFilename: file.original_filename,
    storedFilename: file.stored_filename,
    bucketName: file.bucket_name,
    mimeType: file.mime_type,
    sizeBytes: file.size_bytes,
    uploadedBy: file.uploaded_by,
    tags: file.tags,
    downloadUrl: file.download_url,
    createdAt: file.created_at,
    updatedAt: file.updated_at,
    deletedAt: file.deleted_at,
  };
}

export const filesApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    listFiles: builder.query<FilesListResponse, ListFilesParams | void>({
      query: (params) => {
        const searchParams = new URLSearchParams();
        if (params?.uploadedBy) searchParams.append('uploaded_by', params.uploadedBy);
        if (params?.tags && params.tags.length > 0) {
          searchParams.append('tags', params.tags.join(','));
        }
        if (params?.mimeType) searchParams.append('mime_type', params.mimeType);
        if (params?.bucketName) searchParams.append('bucket_name', params.bucketName);
        if (params?.limit) searchParams.append('limit', params.limit.toString());
        if (params?.offset) searchParams.append('offset', params.offset.toString());
        if (params?.sortColumn) searchParams.append('sort_column', params.sortColumn);
        if (params?.sortDirection) searchParams.append('sort_direction', params.sortDirection);

        const queryString = searchParams.toString();
        return {
          url: `${endpoints.files.base}${queryString ? `?${queryString}` : ''}`,
          method: 'GET',
        };
      },
      transformResponse: (response: unknown) => {
        const data = response as { files: FileResponse[]; total: number };
        return {
          files: data.files.map((file) => transformFile(file)),
          total: data.total,
        };
      },
      providesTags: (result) =>
        result
          ? [
              ...result.files.map(({ id }) => ({ type: 'File' as const, id })),
              { type: 'File', id: 'LIST' },
            ]
          : [{ type: 'File', id: 'LIST' }],
    }),

    getFile: builder.query<File, string>({
      query: (id) => ({
        url: endpoints.files.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const file = response as FileResponse;
        return transformFile(file);
      },
      providesTags: (_result, _error, id) => [{ type: 'File', id }],
    }),

    uploadFile: builder.mutation<File, UploadFileParams>({
      queryFn: async ({ file, bucketName, tags }, api, extraOptions) => {
        console.log('[API] uploadFile queryFn - START (unconditional log):', {
          hasFile: !!file,
          fileName: file instanceof File ? file.name : file instanceof Blob ? 'Blob' : 'unknown',
          bucketName,
          tags,
        });

        const getBaseURL = (): string => {
          const envBaseURL = import.meta.env.VITE_API_BASE_URL;
          if (envBaseURL) {
            return envBaseURL;
          }
          return 'http://asto-lms.local';
        };

        const baseURL = getBaseURL();
        const url = `${baseURL}${endpoints.files.upload}`;

        const formData = new FormData();
        formData.append('file', file as unknown as Blob);
        if (bucketName) {
          formData.append('bucket_name', bucketName);
        }
        if (tags && tags.length > 0) {
          formData.append('tags', tags.join(','));
        }

        const dashboardToken = storage.getDashboardAccessToken();
        const publicToken = storage.getAccessToken();
        const token = dashboardToken || publicToken;
        const bearerToken = token && typeof token === 'string' && token.trim().length > 0
          ? `Bearer ${token.trim()}`
          : null;

        const headers: HeadersInit = {};
        if (bearerToken) {
          headers['Authorization'] = bearerToken;
        }

        console.log('[API] uploadFile queryFn - Making fetch request (unconditional log):', {
          url,
          hasFormData: true,
          formDataKeys: Array.from(formData.keys()),
          hasToken: !!bearerToken,
          headersKeys: Object.keys(headers),
        });

        try {
          const response = await fetch(url, {
            method: 'POST',
            body: formData,
            headers: Object.keys(headers).length > 0 ? headers : undefined,
          });

          console.log('[API] uploadFile queryFn - Response received (unconditional log):', {
            status: response.status,
            statusText: response.statusText,
            ok: response.ok,
            contentType: response.headers.get('Content-Type'),
          });

          if (!response.ok) {
            const clonedResponse = response.clone();
            let errorData: unknown;
            try {
             
              errorData = await response.json();
            } catch {
              try {
                errorData = await clonedResponse.text();
              } catch {
                errorData = `HTTP ${response.status}: ${response.statusText}`;
              }
            }

            console.error('[API] uploadFile queryFn - Error response (unconditional log):', {
              status: response.status,
              errorData,
            });

            return {
              error: {
                status: response.status,
                data: errorData,
              },
            };
          }

          const data = await response.json();
          const fileResponse = data as FileResponse;
          const transformedFile = transformFile(fileResponse);

          console.log('[API] uploadFile queryFn - Success (unconditional log):', {
            fileId: transformedFile.id,
            fileName: transformedFile.originalFilename,
          });

          return { data: transformedFile };
        } catch (error) {
          console.error('[API] uploadFile queryFn - Fetch error (unconditional log):', error);
          return {
            error: {
              status: 'FETCH_ERROR' as const,
              error: error instanceof Error ? error.message : String(error),
            },
          };
        }
      },
      invalidatesTags: [{ type: 'File', id: 'LIST' }],
    }),

    deleteFile: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.files.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'File', id },
        { type: 'File', id: 'LIST' },
      ],
    }),
  }),
});

export const {
  useListFilesQuery,
  useLazyListFilesQuery,
  useGetFileQuery,
  useUploadFileMutation,
  useDeleteFileMutation,
} = filesApi;

