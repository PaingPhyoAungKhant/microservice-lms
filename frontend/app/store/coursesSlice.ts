import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

export type Course = {
  id: string;
  title: string;
  price: number;
  description: string;
  image: string;
};

const initialState: { courses: Course[] } = {
  courses: [
    {
      id: '1',
      title: 'Full-Stack Web Development',
      price: 200,
      description: 'Build dynamic websites and web apps...',
      image: '/courses/cardimg1.jpg',
    },
    {
      id: '2',
      title: 'UI/UX Design Fundamentals',
      price: 150,
      description: 'Learn the basics of user interface and user experience design.',
      image: '/courses/cardimg2.jpg',
    },
  ],
};

const coursesSlice = createSlice({
  name: 'courses',
  initialState,
  reducers: {
    setCourses(state, action: PayloadAction<Course[]>) {
      state.courses = action.payload;
    },
  },
});

export const { setCourses } = coursesSlice.actions;
export default coursesSlice.reducer;
