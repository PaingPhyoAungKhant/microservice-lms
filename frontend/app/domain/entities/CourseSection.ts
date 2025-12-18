export interface CourseSection {
  id: string;
  courseOfferingId: string;
  name: string;
  description: string;
  order: number;
  status: 'draft' | 'published' | 'archived';
  createdAt: string;
  updatedAt: string;
}

