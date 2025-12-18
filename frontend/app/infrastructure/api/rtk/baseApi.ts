import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { storage } from '../../storage/storage';
import { endpoints } from '../endpoints';
import type { ApiException } from '../client';

const getBaseURL = (): string => {
  const envBaseURL = import.meta.env.VITE_API_BASE_URL;
  if (envBaseURL) {
    return envBaseURL;
  }
  return 'http://asto-lms.local';
};

const baseQuery = fetchBaseQuery({
  baseUrl: getBaseURL(),
  prepareHeaders: (headers, { endpoint }) => {
    const dashboardToken = storage.getDashboardAccessToken();
    const publicToken = storage.getAccessToken();
    const token = dashboardToken || publicToken;

    if (import.meta.env.DEV) {
      console.debug('[API] prepareHeaders called', {
        endpoint,
        hasDashboardToken: !!dashboardToken,
        hasPublicToken: !!publicToken,
        tokenLength: token ? token.length : 0,
      });
    }

    if (token && typeof token === 'string' && token.trim().length > 0) {
      const bearerToken = `Bearer ${token.trim()}`;
      headers.set('Authorization', bearerToken);
      
      if (import.meta.env.DEV) {
        const tokenType = dashboardToken ? 'dashboard' : 'public';
        console.debug(`[API] Authorization header set in prepareHeaders (${tokenType} token)`, {
          headerValue: bearerToken.substring(0, 30) + '...',
          headerExists: headers.has('Authorization'),
        });
      }
    } else if (import.meta.env.DEV) {
      console.warn('[API] No valid access token available for request', {
        endpoint,
        hasDashboardToken: !!dashboardToken,
        hasPublicToken: !!publicToken,
        tokenType: typeof token,
        dashboardTokenValue: dashboardToken ? dashboardToken.substring(0, 20) + '...' : null,
        publicTokenValue: publicToken ? publicToken.substring(0, 20) + '...' : null,
      });
    }

    return headers;
  },
  fetchFn: async (input, init) => {
    const isFormData = init?.body instanceof FormData;
    const url = typeof input === 'string' ? input : input?.url || 'unknown';
    
    const dashboardToken = storage.getDashboardAccessToken();
    const publicToken = storage.getAccessToken();
    const token = dashboardToken || publicToken;
    const bearerToken = token && typeof token === 'string' && token.trim().length > 0
      ? `Bearer ${token.trim()}`
      : null;

    if (import.meta.env.DEV) {
      console.debug('[API] fetchFn called', {
        url,
        isFormData,
        hasBody: !!init?.body,
        bodyType: init?.body?.constructor?.name,
        initHeadersType: init?.headers?.constructor?.name,
        initHasHeaders: !!init?.headers,
      });
    }

    if (isFormData) {
      const initHeaders = init?.headers;
      let initHeadersInfo: any = {
        type: initHeaders?.constructor?.name || 'null/undefined',
        hasHeaders: !!initHeaders,
      };
      
      if (initHeaders instanceof Headers) {
        initHeadersInfo.headerKeys = Array.from(initHeaders.keys());
        initHeadersInfo.contentType = initHeaders.get('Content-Type');
        console.log('[API] FormData - RTK Query passed Headers object (unconditional log):', {
          headerKeys: initHeadersInfo.headerKeys,
          contentType: initHeadersInfo.contentType,
        });
      } else if (initHeaders && typeof initHeaders === 'object') {
        initHeadersInfo.headerKeys = Object.keys(initHeaders);
        initHeadersInfo.contentType = (initHeaders as any)['Content-Type'] || (initHeaders as any)['content-type'];
        console.log('[API] FormData - RTK Query passed headers object (unconditional log):', {
          headerKeys: initHeadersInfo.headerKeys,
          contentType: initHeadersInfo.contentType,
          headersObject: initHeaders,
        });
      } else {
        console.log('[API] FormData - RTK Query passed no headers (unconditional log)');
      }
      
      let cleanedHeaders: Headers | Record<string, string> | undefined;
      
      if (initHeaders instanceof Headers) {
        cleanedHeaders = new Headers(initHeaders);
        if (cleanedHeaders.has('Content-Type')) {
          console.log('[API] FormData - Removing Content-Type from Headers object (unconditional log)');
          cleanedHeaders.delete('Content-Type');
        }
        if (cleanedHeaders.has('content-type')) {
          console.log('[API] FormData - Removing content-type (lowercase) from Headers object (unconditional log)');
          cleanedHeaders.delete('content-type');
        }
      } else if (initHeaders && typeof initHeaders === 'object') {
        const headersObj = initHeaders as Record<string, string>;
        cleanedHeaders = { ...headersObj };
        if ('Content-Type' in cleanedHeaders) {
          console.log('[API] FormData - Removing Content-Type from headers object (unconditional log)');
          delete cleanedHeaders['Content-Type'];
        }
        if ('content-type' in cleanedHeaders) {
          console.log('[API] FormData - Removing content-type (lowercase) from headers object (unconditional log)');
          delete cleanedHeaders['content-type'];
        }
      }
      
      const formDataHeaders: Record<string, string> = {};
      if (bearerToken) {
        formDataHeaders['Authorization'] = bearerToken;
      }
      
      if (cleanedHeaders instanceof Headers) {
        if (bearerToken) {
          cleanedHeaders.set('Authorization', bearerToken);
        }
        const { headers: _excludedHeaders, ...initWithoutHeaders } = init || {};
        const fetchInit: RequestInit = {
          ...initWithoutHeaders,
          headers: cleanedHeaders,
        };
        
        console.log('[API] FormData - Final fetchInit with cleaned Headers (unconditional log):', {
          hasHeaders: true,
          headerKeys: Array.from(cleanedHeaders.keys()),
          hasContentType: cleanedHeaders.has('Content-Type') || cleanedHeaders.has('content-type'),
        });
        
        return fetch(input, fetchInit);
      } else if (cleanedHeaders && typeof cleanedHeaders === 'object') {
        const finalHeaders = { ...cleanedHeaders, ...formDataHeaders };
        const { headers: _excludedHeaders, ...initWithoutHeaders } = init || {};
        const fetchInit: RequestInit = {
          ...initWithoutHeaders,
          headers: finalHeaders,
        };
        
        console.log('[API] FormData - Final fetchInit with cleaned headers object (unconditional log):', {
          hasHeaders: true,
          headerKeys: Object.keys(finalHeaders),
          hasContentType: 'Content-Type' in finalHeaders || 'content-type' in finalHeaders,
        });
        
        return fetch(input, fetchInit);
      } else {
        const { headers: _excludedHeaders, ...initWithoutHeaders } = init || {};
        const fetchInit: RequestInit = {
          ...initWithoutHeaders,
          ...(Object.keys(formDataHeaders).length > 0 ? { headers: formDataHeaders } : {}),
        };
        
        console.log('[API] FormData - Final fetchInit with minimal headers (unconditional log):', {
          hasHeaders: 'headers' in fetchInit,
          headerKeys: fetchInit.headers ? Object.keys(fetchInit.headers as Record<string, string>) : [],
          hasContentType: fetchInit.headers && (
            (fetchInit.headers instanceof Headers && (fetchInit.headers.has('Content-Type') || fetchInit.headers.has('content-type'))) ||
            (typeof fetchInit.headers === 'object' && ('Content-Type' in fetchInit.headers || 'content-type' in fetchInit.headers))
          ),
        });
        
        return fetch(input, fetchInit);
      }
    } else {
      let headers: Headers;
      if (init?.headers instanceof Headers) {
        headers = new Headers(init.headers);
      } else if (init?.headers) {
        headers = new Headers(init.headers as HeadersInit);
      } else {
        headers = new Headers();
      }
      
      if (!headers.has('Authorization') && bearerToken) {
        headers.set('Authorization', bearerToken);
        
        if (import.meta.env.DEV) {
          const tokenType = dashboardToken ? 'dashboard' : 'public';
          console.debug(`[API] Authorization header added in fetchFn fallback (${tokenType} token)`, {
            headerValue: bearerToken.substring(0, 30) + '...',
            url,
          });
        }
      } else if (import.meta.env.DEV && !bearerToken) {
        if (!url.includes('/auth/')) {
          console.warn('[API] Authorization header MISSING in fetchFn and no valid token available!', {
            url,
            hasDashboardToken: !!dashboardToken,
            hasPublicToken: !!publicToken,
            tokenType: typeof token,
          });
        }
      }
      
      if (!headers.has('Content-Type') && init?.body) {
        headers.set('Content-Type', 'application/json');
        
        if (import.meta.env.DEV) {
          console.debug('[API] Content-Type set to application/json for JSON request', {
            url,
            hasBody: !!init?.body,
          });
        }
      }
      
      if (import.meta.env.DEV) {
        const authHeader = headers.get('Authorization');
        const contentType = headers.get('Content-Type');
        console.debug('[API] JSON request headers', {
          url,
          hasAuthorization: !!authHeader,
          contentType: contentType || 'not set',
          allHeaders: Array.from(headers.keys()),
        });
      }

      return fetch(input, {
        ...init,
        headers,
      });
    }
  },
});

let isRefreshing = false;
let refreshPromise: Promise<boolean> | null = null;

const baseQueryWithReauth = async (args: any, api: any, extraOptions: any) => {
  let result = await baseQuery(args, api, extraOptions);

  if (result.error && result.error.status === 401) {
    const url = typeof args === 'string' ? args : args?.url || '';
    if (url.includes('/auth/refresh-token')) {
      storage.removeAccessToken();
      storage.removeRefreshToken();
      storage.removeDashboardAccessToken();
      storage.removeDashboardRefreshToken();
      return result;
    }

    if (isRefreshing && refreshPromise) {
      const refreshSuccess = await refreshPromise;
      if (refreshSuccess) {
        result = await baseQuery(args, api, extraOptions);
      }
      return result;
    }

    const dashboardToken = storage.getDashboardAccessToken();
    const publicToken = storage.getAccessToken();
    const tokenType = dashboardToken ? 'dashboard' : publicToken ? 'public' : null;

    if (!tokenType) {
      return result;
    }

    const refreshToken = tokenType === 'dashboard' 
      ? storage.getDashboardRefreshToken()
      : storage.getRefreshToken();

    if (!refreshToken) {
      if (tokenType === 'dashboard') {
        storage.clearDashboard();
      } else {
        storage.removeAccessToken();
        storage.removeRefreshToken();
        storage.removeUser();
      }
      return result;
    }

    isRefreshing = true;
    refreshPromise = (async () => {
      try {
        const baseUrl = getBaseURL();
        const refreshUrl = `${baseUrl}${endpoints.auth.refresh}`;
        const refreshResponse = await fetch(refreshUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });

        if (refreshResponse.ok) {
          const refreshData = await refreshResponse.json();
          const newAccessToken = refreshData.access_token;
          const newRefreshToken = refreshData.refresh_token;

          if (!newAccessToken) {
            console.error('Token refresh failed: No access_token in response', refreshData);
            if (tokenType === 'dashboard') {
              storage.clearDashboard();
            } else {
              storage.removeAccessToken();
              storage.removeRefreshToken();
              storage.removeUser();
            }
            return false;
          }

          if (tokenType === 'dashboard') {
            storage.setDashboardAccessToken(newAccessToken);
            if (newRefreshToken) {
              storage.setDashboardRefreshToken(newRefreshToken);
            }
          } else {
            storage.setAccessToken(newAccessToken);
            if (newRefreshToken) {
              storage.setRefreshToken(newRefreshToken);
            }
          }

          return true;
        } else {
          const errorData = await refreshResponse.json().catch(() => ({}));
          console.error('Token refresh failed:', refreshResponse.status, errorData);
          
          if (tokenType === 'dashboard') {
            storage.clearDashboard();
          } else {
            storage.removeAccessToken();
            storage.removeRefreshToken();
            storage.removeUser();
          }
          return false;
        }
      } catch (error) {

        console.error('Token refresh request failed:', error);
        if (tokenType === 'dashboard') {
          storage.clearDashboard();
        } else {
          storage.removeAccessToken();
          storage.removeRefreshToken();
          storage.removeUser();
        }
        return false;
      } finally {
        isRefreshing = false;
        refreshPromise = null;
      }
    })();

    const refreshSuccess = await refreshPromise;
    
    if (refreshSuccess) {
      result = await baseQuery(args, api, extraOptions);
    }
  }

  return result;
};

export const baseApi = createApi({
  reducerPath: 'api',
  baseQuery: baseQueryWithReauth,
  tagTypes: ['Auth', 'User', 'Course', 'Enrollment', 'Category', 'File', 'CourseOffering'],
  endpoints: () => ({}),
  refetchOnFocus: false,
  refetchOnReconnect: false,
});

