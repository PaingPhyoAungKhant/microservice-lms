export interface File {
  id: string;
  originalFilename: string;
  storedFilename: string;
  bucketName: string;
  mimeType: string;
  sizeBytes: number;
  uploadedBy: string;
  tags: string[];
  downloadUrl?: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string | null;
}

