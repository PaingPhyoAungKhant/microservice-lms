import { baseApi } from './baseApi';
import { endpoints } from '../endpoints';
import {
  sectionModuleResponseSchema,
  sectionModulesResponseSchema,
  createSectionModuleRequestSchema,
  updateSectionModuleRequestSchema,
  reorderSectionModulesRequestSchema,
} from '../../validation/schemas';
import type { SectionModule } from '../../../domain/entities/SectionModule';
import type { z } from 'zod';

function transformSectionModule(module: z.infer<typeof sectionModuleResponseSchema>): SectionModule {
  return {
    id: module.id,
    courseSectionId: module.course_section_id,
    contentId: module.content_id ?? undefined,
    name: module.name,
    description: module.description,
    contentType: module.content_type,
    contentStatus: module.content_status,
    order: module.order,
    createdAt: module.created_at,
    updatedAt: module.updated_at,
  };
}

interface GetSectionModulesParams {
  sectionId: string;
}

interface SectionModulesListResponse {
  modules: SectionModule[];
}

interface CreateSectionModuleRequest {
  name: string;
  description?: string;
  contentType: 'zoom';
  order?: number;
}

interface UpdateSectionModuleRequest {
  name?: string;
  description?: string;
  order?: number;
}

interface ReorderItem {
  id: string;
  order: number;
}

function transformCreateModuleRequest(data: CreateSectionModuleRequest): z.infer<typeof createSectionModuleRequestSchema> {
  return {
    name: data.name,
    description: data.description,
    content_type: data.contentType,
    order: data.order,
  };
}

export const sectionModulesApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getSectionModules: builder.query<SectionModulesListResponse, GetSectionModulesParams>({
      query: ({ sectionId }) => ({
        url: endpoints.sectionModules.bySection(sectionId),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = sectionModulesResponseSchema.parse(response);
        return {
          modules: validated.modules.map(transformSectionModule),
        };
      },
      providesTags: (result, _error, { sectionId }) =>
        result
          ? [
              ...result.modules.map(({ id }) => ({ type: 'SectionModule' as const, id })),
              { type: 'SectionModule', id: `LIST-${sectionId}` },
            ]
          : [{ type: 'SectionModule', id: `LIST-${sectionId}` }],
    }),

    getSectionModule: builder.query<SectionModule, string>({
      query: (id) => ({
        url: endpoints.sectionModules.byId(id),
        method: 'GET',
      }),
      transformResponse: (response: unknown) => {
        const validated = sectionModuleResponseSchema.parse(response);
        return transformSectionModule(validated);
      },
      providesTags: (_result, _error, id) => [{ type: 'SectionModule', id }],
    }),

    createSectionModule: builder.mutation<SectionModule, { sectionId: string; data: CreateSectionModuleRequest }>({
      query: ({ sectionId, data }) => {
        
        const transformedData = transformCreateModuleRequest(data);
        return {
          url: endpoints.sectionModules.bySection(sectionId),
          method: 'POST',
          body: createSectionModuleRequestSchema.parse(transformedData),
        };
      },
      transformResponse: (response: unknown) => {
        const validated = sectionModuleResponseSchema.parse(response);
        return transformSectionModule(validated);
      },
      invalidatesTags: (_result, _error, { sectionId }) => [
        { type: 'SectionModule', id: 'LIST' },
        { type: 'SectionModule', id: `LIST-${sectionId}` },
      ],
    }),

    updateSectionModule: builder.mutation<SectionModule, { id: string; data: UpdateSectionModuleRequest }>({
      query: ({ id, data }) => ({
        url: endpoints.sectionModules.byId(id),
        method: 'PUT',
        body: updateSectionModuleRequestSchema.parse(data),
      }),
      transformResponse: (response: unknown) => {
        const validated = sectionModuleResponseSchema.parse(response);
        return transformSectionModule(validated);
      },
      invalidatesTags: (_result, _error, { id }) => [
        { type: 'SectionModule', id },
        { type: 'SectionModule', id: 'LIST' },
      ],
    }),

    deleteSectionModule: builder.mutation<void, string>({
      query: (id) => ({
        url: endpoints.sectionModules.byId(id),
        method: 'DELETE',
      }),
      invalidatesTags: (_result, _error, id) => [
        { type: 'SectionModule', id },
        { type: 'SectionModule', id: 'LIST' },
      ],
    }),

    reorderSectionModules: builder.mutation<void, { sectionId: string; items: ReorderItem[] }>({
      query: ({ sectionId, items }) => ({
        url: endpoints.sectionModules.reorder(sectionId),
        method: 'PUT',
        body: reorderSectionModulesRequestSchema.parse({ items }),
      }),
      invalidatesTags: (_result, _error, { sectionId }) => [
        { type: 'SectionModule', id: 'LIST' },
        { type: 'SectionModule', id: `LIST-${sectionId}` },
      ],
    }),
  }),
});

export const {
  useGetSectionModulesQuery,
  useLazyGetSectionModulesQuery,
  useGetSectionModuleQuery,
  useCreateSectionModuleMutation,
  useUpdateSectionModuleMutation,
  useDeleteSectionModuleMutation,
  useReorderSectionModulesMutation,
} = sectionModulesApi;

