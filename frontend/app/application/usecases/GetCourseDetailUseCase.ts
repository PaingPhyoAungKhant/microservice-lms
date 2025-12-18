import type { Course } from '../../domain/entities/Course';
import type { ICourseRepository } from '../../domain/repositories/ICourseRepository';

export interface GetCourseDetailInput {
  courseId: string;
}

export interface GetCourseDetailOutput {
  course: Course | null;
}

export class GetCourseDetailUseCase {
  constructor(private courseRepository: ICourseRepository) {}

  async execute(input: GetCourseDetailInput): Promise<GetCourseDetailOutput> {
    const course = await this.courseRepository.findById(input.courseId);

    return {
      course,
    };
  }
}

