import type { AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios';

export interface ApiError {
  message: string;
  status: number;
  code?: string;
  errors?: Record<string, string[]>;
}

export class ApiException extends Error {
  constructor(
    public message: string,
    public status: number,
    public code?: string,
    public errors?: Record<string, string[]>
  ) {
    super(message);
    this.name = 'ApiException';
  }
}

import { storage } from '../storage/storage';

export const requestInterceptor = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
  const dashboardToken = storage.getDashboardAccessToken();
  const publicToken = storage.getAccessToken();
  const token = dashboardToken || publicToken;
  
  if (token && typeof token === 'string' && token.trim().length > 0 && config.headers) {
    config.headers.Authorization = `Bearer ${token.trim()}`;
    
    if (import.meta.env.DEV) {
      const tokenType = dashboardToken ? 'dashboard' : 'public';
      console.debug(`[Axios Interceptor] Using ${tokenType} access token for request`);
    }
  } else if (import.meta.env.DEV && config.url && !config.url.includes('/auth/')) {
    console.warn('[Axios Interceptor] No valid access token available for request', {
      url: config.url,
      hasDashboardToken: !!dashboardToken,
      hasPublicToken: !!publicToken,
      tokenType: typeof token,
    });
  }

  if (config.headers) {
    config.headers['Content-Type'] = config.headers['Content-Type'] || 'application/json';
  }

  return config;
};

export const responseInterceptor = (response: AxiosResponse): AxiosResponse => {
  return response;
};

export const errorInterceptor = (error: any): Promise<ApiException> => {
  if (error.response) {
    const { status, data } = error.response;
    const message = data?.message || data?.error || 'An error occurred';
    const code = data?.code;
    const errors = data?.errors;

    return Promise.reject(new ApiException(message, status, code, errors));
  } else if (error.request) {
    return Promise.reject(new ApiException('Network error. Please check your connection.', 0));
  } else {
    return Promise.reject(new ApiException(error.message || 'An unexpected error occurred', 0));
  }
};

