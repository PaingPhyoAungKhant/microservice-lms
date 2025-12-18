import type { Category } from './Category';

export interface Course {
  id: string;
  name: string;
  description: string;
  thumbnailId?: string | null;
  thumbnailUrl?: string;
  categories?: Category[];
  createdAt: string;
  updatedAt: string;
}

