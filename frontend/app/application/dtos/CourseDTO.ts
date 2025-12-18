import type { Course } from '../../domain/entities/Course';

export interface CourseDTO extends Course {}

export interface CourseListDTO {
  courses: CourseDTO[];
  total: number;
}

export interface CreateCourseDTO {
  title: string;
  description: string;
  price: number;
  image?: string;
  category?: string;
}

export interface UpdateCourseDTO {
  title?: string;
  description?: string;
  price?: number;
  image?: string;
  category?: string;
}

