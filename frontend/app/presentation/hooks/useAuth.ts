import { useCallback, useEffect } from 'react';
import { useNavigate } from 'react-router';
import { useDispatch, useSelector } from 'react-redux';
import {
  useLoginMutation,
  useRegisterMutation,
  useVerifyTokenQuery,
} from '../../infrastructure/api/rtk/authApi';
import { clearPublicAuth, setPublicUser, setPublicLoading } from '../../infrastructure/store/authSlice';
import { storage } from '../../infrastructure/storage/storage';
import { ROUTES } from '../../shared/constants/routes';
import type { RootState } from '../../store';
import type { User } from '../../domain/entities/User';
import type { FetchBaseQueryError } from '@reduxjs/toolkit/query';

interface ApiException {
  message: string;
  status: number;
  code?: string;
  errors?: Record<string, string[]>;
}

export interface UseAuthReturn {
  user: User | null;
  loading: boolean;
  error: ApiException | null;
  login: (email: string, password: string) => Promise<void>;
  register: (username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  reset: () => void;
}

export function useAuth(): UseAuthReturn {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const publicUser = useSelector((state: RootState) => state.auth.publicUser);
  const publicIsLoading = useSelector((state: RootState) => state.auth.publicIsLoading);

  const [loginMutation, { isLoading: loginLoading, error: loginError }] = useLoginMutation();
  const [registerMutation, { isLoading: registerLoading, error: registerError }] = useRegisterMutation();
  const accessToken = storage.getAccessToken();
  const { isLoading: verifyLoading, isError: verifyError } = useVerifyTokenQuery(undefined, {
    skip: !accessToken,
    refetchOnMountOrArgChange: false,
    refetchOnFocus: false,
    refetchOnReconnect: false,
  });

  useEffect(() => {
    if (verifyError && accessToken) {
      storage.removeAccessToken();
      storage.removeRefreshToken();
      storage.removeUser();
      dispatch(clearPublicAuth());
    }
  }, [verifyError, accessToken, dispatch]);

  const loading = loginLoading || registerLoading || verifyLoading || publicIsLoading;

  const getError = (): ApiException | null => {
    if (loginError) {
      const err = loginError as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Login failed',
        status: err.status as number || 500,
      };
    }
    if (registerError) {
      const err = registerError as FetchBaseQueryError;
      return {
        message: (err.data as any)?.message || 'Registration failed',
        status: err.status as number || 500,
      };
    }
    return null;
  };

  const login = useCallback(
    async (email: string, password: string) => {
      dispatch(setPublicLoading(true));
      try {
        await loginMutation({ email, password }).unwrap();
      } catch (err) {
        dispatch(setPublicLoading(false));
        throw err;
      } finally {
        dispatch(setPublicLoading(false));
      }
    },
    [loginMutation, dispatch]
  );

  const register = useCallback(
    async (username: string, email: string, password: string) => {
      dispatch(setPublicLoading(true));
      try {
        await registerMutation({ username, email, password }).unwrap();
      } catch (err) {
        dispatch(setPublicLoading(false));
        throw err;
      } finally {
        dispatch(setPublicLoading(false));
      }
    },
    [registerMutation, dispatch]
  );

  const logout = useCallback(() => {
    dispatch(clearPublicAuth());
    navigate(ROUTES.HOME);
  }, [dispatch, navigate]);

  const reset = useCallback(() => {
  }, []);

  return {
    user: publicUser,
    loading,
    error: getError(),
    login,
    register,
    logout,
    reset,
  };
}
