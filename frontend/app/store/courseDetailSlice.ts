import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

type CourseDetail = {
  id: string;
  title: string;
  price: number;
  description: string;
  image: string;
};

const initialState: { course: CourseDetail | null } = {
  course: {
    id: '1',
    title: 'Full-Stack Web Development',
    price: 200,
    description:
      'Build dynamic websites and web apps using HTML, CSS, Javascript, React, Node.js and MongoDB.',
    image: '/courses/cardimg1.jpg',
  },
};

const courseDetailSlice = createSlice({
  name: 'courseDetail',
  initialState,
  reducers: {
    setCourseDetail(state, action: PayloadAction<CourseDetail>) {
      state.course = action.payload;
    },
  },
});

export const { setCourseDetail } = courseDetailSlice.actions;
export default courseDetailSlice.reducer;
