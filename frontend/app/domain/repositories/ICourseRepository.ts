import type { Course } from '../entities/Course';

export interface CourseQuery {
  searchQuery?: string;
  category?: string;
  instructorId?: string;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

export interface CourseListResult {
  courses: Course[];
  total: number;
}

export interface ICourseRepository {
  findById(id: string): Promise<Course | null>;
  find(query: CourseQuery): Promise<CourseListResult>;
  create(course: Omit<Course, 'id' | 'createdAt' | 'updatedAt'>): Promise<Course>;
  update(id: string, course: Partial<Course>): Promise<Course>;
  delete(id: string): Promise<void>;
}

