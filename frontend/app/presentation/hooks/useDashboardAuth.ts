import { useCallback } from 'react';
import { useNavigate } from 'react-router';
import { useDispatch, useSelector } from 'react-redux';
import { useLoginDashboardMutation } from '../../infrastructure/api/rtk/authApi';
import { clearDashboardAuth, setDashboardUser, setDashboardLoading } from '../../infrastructure/store/authSlice';
import { storage } from '../../infrastructure/storage/storage';
import type { RootState } from '../../store';
import type { User } from '../../domain/entities/User';
import type { FetchBaseQueryError } from '@reduxjs/toolkit/query';
import { ROUTES } from '../../shared/constants/routes';

interface ApiException {
  message: string;
  status: number;
  code?: string;
  errors?: Record<string, string[]>;
}

export interface UseDashboardAuthReturn {
  user: User | null;
  loading: boolean;
  error: ApiException | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  reset: () => void;
}

export function useDashboardAuth(): UseDashboardAuthReturn {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const dashboardUser = useSelector((state: RootState) => state.auth.dashboardUser);
  const dashboardIsLoading = useSelector((state: RootState) => state.auth.dashboardIsLoading);

  const [loginDashboardMutation, { isLoading: loginLoading, error: loginError }] = useLoginDashboardMutation();

  const loading = loginLoading || dashboardIsLoading;

  const getError = (): ApiException | null => {
    if (loginError) {
      const err = loginError as FetchBaseQueryError;
      let message = 'Login failed';
      
      if (err.status === 404) {
        message = 'Login endpoint not found. Please check if the auth service is running and the API gateway is configured correctly.';
      } else if (err.status === 'PARSING_ERROR') {
        message = 'Invalid response from server. The login endpoint may not be available.';
      } else if (err.data) {
        
        if (typeof err.data === 'object' && 'message' in err.data) {
          message = (err.data as any).message;
        } else if (typeof err.data === 'string') {
          message = err.data;
        }
      }
      
      return {
        message,
        status: (err.status as number) || 500,
      };
    }
    return null;
  };

  const login = useCallback(
    async (email: string, password: string) => {
      dispatch(setDashboardLoading(true));
      try {
        const result = await loginDashboardMutation({ email, password }).unwrap();
        
        const role = result.user.role;
        if (role === 'admin') {
          navigate(ROUTES.ADMIN_USERS);
        } else if (role === 'instructor') {
          navigate(ROUTES.INSTRUCTOR_DASHBOARD);
        } else if (role === 'student') {
          navigate(ROUTES.STUDENT_DASHBOARD);
        } else {
          navigate(ROUTES.DASHBOARD);
        }
      } catch (err) {
        dispatch(setDashboardLoading(false));
        throw err;
      } finally {
        dispatch(setDashboardLoading(false));
      }
    },
    [loginDashboardMutation, dispatch, navigate]
  );

  const logout = useCallback(() => {
    dispatch(clearDashboardAuth());
    navigate(ROUTES.DASHBOARD);
  }, [dispatch, navigate]);

  const reset = useCallback(() => {
  }, []);

  return {
    user: dashboardUser,
    loading,
    error: getError(),
    login,
    logout,
    reset,
  };
}
