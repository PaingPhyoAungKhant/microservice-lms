import type {
  ICourseRepository,
  CourseQuery,
  CourseListResult,
} from '../../domain/repositories/ICourseRepository';
import type { Course } from '../../domain/entities/Course';
import { apiClient } from '../api/client';
import { endpoints } from '../api/endpoints';
import {
  courseResponseSchema,
  coursesResponseSchema,
  type Course as CourseType,
} from '../validation/schemas';
import { ApiException } from '../api/client';

export class CourseRepository implements ICourseRepository {
  async findById(id: string): Promise<Course | null> {
    try {
      const data = await apiClient.get<CourseType>(endpoints.courses.byId(id));
      const validated = courseResponseSchema.parse(data);
      return this.mapToDomain(validated);
    } catch (error) {
      if (error instanceof ApiException && error.status === 404) {
        return null;
      }
      throw error;
    }
  }

  async find(query: CourseQuery): Promise<CourseListResult> {
    try {
      const params = new URLSearchParams();
      if (query.searchQuery) params.append('search', query.searchQuery);
      if (query.category) params.append('category', query.category);
      if (query.instructorId) params.append('instructorId', query.instructorId);
      if (query.limit) params.append('limit', query.limit.toString());
      if (query.offset) params.append('offset', query.offset.toString());
      if (query.sortColumn) params.append('sortColumn', query.sortColumn);
      if (query.sortDirection) params.append('sortDirection', query.sortDirection);

      const url = `${endpoints.courses.find}?${params.toString()}`;
      const data = await apiClient.get<CourseType[]>(url);
      const validated = coursesResponseSchema.parse(data);

      return {
        courses: validated.map((course) => this.mapToDomain(course)),
        total: validated.length,
      };
    } catch (error) {
      throw error;
    }
  }

  async create(course: Omit<Course, 'id' | 'createdAt' | 'updatedAt'>): Promise<Course> {
    try {
      const data = await apiClient.post<CourseType>(endpoints.courses.base, course);
      const validated = courseResponseSchema.parse(data);
      return this.mapToDomain(validated);
    } catch (error) {
      throw error;
    }
  }

  async update(id: string, course: Partial<Course>): Promise<Course> {
    try {
      const data = await apiClient.put<CourseType>(endpoints.courses.byId(id), course);
      const validated = courseResponseSchema.parse(data);
      return this.mapToDomain(validated);
    } catch (error) {
      throw error;
    }
  }

  async delete(id: string): Promise<void> {
    try {
      await apiClient.delete<void>(endpoints.courses.byId(id));
    } catch (error) {
      throw error;
    }
  }

  private mapToDomain(course: CourseType): Course {
    return {
      id: course.id,
      name: course.name,
      description: course.description,
      thumbnailId: course.thumbnail_id ?? undefined,
      thumbnailUrl: course.thumbnail_url ?? undefined,
      categories: course.categories?.map((category) => ({
        id: category.id,
        name: category.name,
        description: category.description,
        createdAt: category.created_at,
        updatedAt: category.updated_at,
      })) ?? [],
      createdAt: course.created_at,
      updatedAt: course.updated_at,
    };
  }
}

