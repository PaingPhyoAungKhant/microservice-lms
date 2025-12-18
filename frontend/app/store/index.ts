import { configureStore } from '@reduxjs/toolkit';
import courseDetailReducer from './courseDetailSlice';
import coursesReducer from './coursesSlice';

export const store = configureStore({
  reducer: {
    courseDetail: courseDetailReducer,
    courses: coursesReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
