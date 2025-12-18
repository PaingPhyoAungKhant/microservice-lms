export interface SectionModule {
  id: string;
  courseSectionId: string;
  contentId?: string;
  name: string;
  description: string;
  contentType: 'zoom';
  contentStatus: 'draft' | 'pending' | 'created';
  order: number;
  createdAt: string;
  updatedAt: string;
}

