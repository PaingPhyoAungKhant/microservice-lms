import { endpoints } from './endpoints';
import { storage } from '../storage/storage';


const getBaseURL = (): string => {
  const envBaseURL = import.meta.env.VITE_API_BASE_URL;
  if (envBaseURL) {
    return envBaseURL;
  }
  return 'http://asto-lms.local';
};


export function getFileDownloadUrl(bucket: string, fileId: string): string {
  const baseURL = getBaseURL();
  const downloadPath = endpoints.buckets.download(bucket, fileId);
  return `${baseURL}${downloadPath}`;
}


export function getErrorMessage(error: unknown, defaultMessage: string = 'An error occurred'): string {
  if (!error) {
    return defaultMessage;
  }


  if (typeof error === 'object' && 'status' in error && 'data' in error) {
    const fetchError = error as { status: number | string; data?: unknown };
    
    if (fetchError.data && typeof fetchError.data === 'object') {
      const errorData = fetchError.data as Record<string, unknown>;
      
      if (typeof errorData.message === 'string' && errorData.message) {
        return errorData.message;
      }
      
      if (typeof errorData.error === 'string' && errorData.error) {
        return errorData.error;
      }
    }
    
    if (fetchError.status === 'FETCH_ERROR') {
      return 'Network error. Please check your connection.';
    }
    
    if (fetchError.status === 'PARSING_ERROR') {
      return 'Failed to parse server response.';
    }
  }

  if (error instanceof Error) {
    return error.message || defaultMessage;
  }

  if (typeof error === 'string') {
    return error;
  }

  return defaultMessage;
}

export async function downloadFileWithAuth(
  bucket: string,
  fileId: string,
  filename?: string
): Promise<void> {
  const downloadUrl = getFileDownloadUrl(bucket, fileId);

  const dashboardToken = storage.getDashboardAccessToken();
  const publicToken = storage.getAccessToken();
  const token = dashboardToken || publicToken;

  if (!token) {
    throw new Error('No authentication token available. Please log in again.');
  }

  const bearerToken = `Bearer ${token.trim()}`;

  const response = await fetch(downloadUrl, {
    method: 'GET',
    headers: {
      'Authorization': bearerToken,
      'Accept': 'application/octet-stream',
    },
  });

  if (!response.ok) {
    let errorMessage = `Download failed: ${response.status} ${response.statusText}`;

    try {
      const errorData = await response.json();
      if (errorData && typeof errorData === 'object') {
        if (typeof errorData.message === 'string' && errorData.message) {
          errorMessage = errorData.message;
        } else if (typeof errorData.error === 'string' && errorData.error) {
          errorMessage = errorData.error;
        }
      }
    } catch {
      const text = await response.text().catch(() => '');
      if (text) {
        errorMessage = text;
      }
    }

    throw new Error(errorMessage);
  }

  let downloadFilename = filename;
  if (!downloadFilename) {
    const contentDisposition = response.headers.get('Content-Disposition');
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
      if (filenameMatch && filenameMatch[1]) {
        downloadFilename = filenameMatch[1].replace(/['"]/g, '');
      }
    }
  }

  if (!downloadFilename) {
    downloadFilename = `download-${fileId}`;
  }

  const blob = await response.blob();

  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = downloadFilename;
  document.body.appendChild(link);
  link.click();

  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

