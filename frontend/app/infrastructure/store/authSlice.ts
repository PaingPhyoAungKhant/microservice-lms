import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { User } from '../../domain/entities/User';
import { storage } from '../storage/storage';

interface AuthState {
  // Public auth
  publicUser: User | null;
  publicIsAuthenticated: boolean;
  publicIsLoading: boolean;

  // Dashboard auth
  dashboardUser: User | null;
  dashboardIsAuthenticated: boolean;
  dashboardIsLoading: boolean;
}

const getInitialState = (): AuthState => {

  if (typeof window === 'undefined') {
    return {
      publicUser: null,
      publicIsAuthenticated: false,
      publicIsLoading: false,
      dashboardUser: null,
      dashboardIsAuthenticated: false,
      dashboardIsLoading: false,
    };
  }

  return {
    publicUser: storage.getUser(),
    publicIsAuthenticated: !!storage.getAccessToken(),
    publicIsLoading: false,
    dashboardUser: storage.getDashboardUser(),
    dashboardIsAuthenticated: !!storage.getDashboardAccessToken(),
    dashboardIsLoading: false,
  };
};

const authSlice = createSlice({
  name: 'auth',
  initialState: getInitialState,
  reducers: {
    // Public auth actions
    setPublicUser: (state, action: PayloadAction<User | null>) => {
      state.publicUser = action.payload;
      state.publicIsAuthenticated = !!action.payload;
    },
    setPublicLoading: (state, action: PayloadAction<boolean>) => {
      state.publicIsLoading = action.payload;
    },
    clearPublicAuth: (state) => {
      state.publicUser = null;
      state.publicIsAuthenticated = false;
      storage.removeAccessToken();
      storage.removeRefreshToken();
      storage.removeUser();
    },

    // Dashboard auth actions
    setDashboardUser: (state, action: PayloadAction<User | null>) => {
      state.dashboardUser = action.payload;
      state.dashboardIsAuthenticated = !!action.payload;
    },
    setDashboardLoading: (state, action: PayloadAction<boolean>) => {
      state.dashboardIsLoading = action.payload;
    },
    clearDashboardAuth: (state) => {
      state.dashboardUser = null;
      state.dashboardIsAuthenticated = false;
      storage.clearDashboard();
    },
  },
});

export const {
  setPublicUser,
  setPublicLoading,
  clearPublicAuth,
  setDashboardUser,
  setDashboardLoading,
  clearDashboardAuth,
} = authSlice.actions;

export default authSlice.reducer;

