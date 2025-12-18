import type { ICourseRepository, CourseQuery, CourseListResult } from '../../domain/repositories/ICourseRepository';

export interface GetCoursesInput {
  searchQuery?: string;
  category?: string;
  instructorId?: string;
  limit?: number;
  offset?: number;
  sortColumn?: string;
  sortDirection?: 'asc' | 'desc';
}

export interface GetCoursesOutput {
  courses: CourseListResult['courses'];
  total: number;
}

export class GetCoursesUseCase {
  constructor(private courseRepository: ICourseRepository) {}

  async execute(input: GetCoursesInput): Promise<GetCoursesOutput> {
    const query: CourseQuery = {
      searchQuery: input.searchQuery,
      category: input.category,
      instructorId: input.instructorId,
      limit: input.limit,
      offset: input.offset,
      sortColumn: input.sortColumn,
      sortDirection: input.sortDirection,
    };

    const result = await this.courseRepository.find(query);

    return {
      courses: result.courses,
      total: result.total,
    };
  }
}

