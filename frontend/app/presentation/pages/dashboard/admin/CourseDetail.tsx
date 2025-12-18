import React, { useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router';
import {
  Box,
  Typography,
  Paper,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import {
  useGetCourseQuery,
  useCreateCourseMutation,
  useUpdateCourseMutation,
} from '../../../../infrastructure/api/rtk/coursesApi';
import { useDashboardAuth } from '../../../hooks/useDashboardAuth';
import DashboardLayout from '../../../components/dashboard/Layout';
import CourseForm from '../../../components/course/CourseForm';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import Success from '../../../components/common/Success';
import Button from '../../../components/common/Button';
import IconButton from '@mui/material/IconButton';
import { ROUTES } from '../../../../shared/constants/routes';

export default function CourseDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user: currentUser, logout } = useDashboardAuth();
  const isNew = id === 'new';

  const {
    data: course,
    isLoading: loadingCourse,
    error: courseError,
  } = useGetCourseQuery(id!, {
    skip: isNew,
  });

  const [createCourse, { isLoading: creating, error: createError }] = useCreateCourseMutation();
  const [updateCourse, { isLoading: updating, error: updateError }] = useUpdateCourseMutation();
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleSubmit = useCallback(
    async (data: {
      name: string;
      description: string;
      thumbnailId?: string | null;
      categoryIds: string[];
    }) => {
      try {
        if (isNew) {
          await createCourse({
            name: data.name,
            description: data.description,
            thumbnailId: data.thumbnailId,
            categoryIds: data.categoryIds,
          }).unwrap();
          setSuccessMessage('Course created successfully');
          setTimeout(() => {
            navigate(ROUTES.ADMIN_COURSES);
          }, 1500);
        } else {
          await updateCourse({
            id: id!,
            data: {
              name: data.name,
              description: data.description,
              thumbnailId: data.thumbnailId,
              categoryIds: data.categoryIds,
            },
          }).unwrap();
          setSuccessMessage('Course updated successfully');
        }
      } catch (err) {
        console.error('Save failed:', err);
      }
    },
    [isNew, id, createCourse, updateCourse, navigate]
  );

  const handleCancel = useCallback(() => {
    navigate(ROUTES.ADMIN_COURSES);
  }, [navigate]);

  if (!currentUser || currentUser.role !== 'admin') {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h5" color="error">
          Access Denied. Admin role required.
        </Typography>
      </Box>
    );
  }

  const error = createError || updateError || (isNew ? null : courseError);
  const loading = creating || updating || (!isNew && loadingCourse);

  return (
    <DashboardLayout user={currentUser} onLogout={logout}>
      <Box>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
          <IconButton onClick={handleCancel} sx={{ mr: 2 }}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h4" component="h1" sx={{ fontWeight: 600 }}>
            {isNew ? 'Create Course' : 'Edit Course'}
          </Typography>
        </Box>

        {successMessage && (
          <Success
            message={successMessage}
            autoDismiss
            autoDismissDelay={3000}
            onDismiss={() => setSuccessMessage(null)}
            fullWidth
          />
        )}

        {!isNew && loadingCourse && <Loading variant="skeleton" fullWidth />}

        {!isNew && courseError && (
          <Error
            message="Failed to load course"
            onRetry={() => window.location.reload()}
            fullWidth
          />
        )}

        {(!isNew ? course : true) && (
          <CourseForm
            course={isNew ? null : course || null}
            onSubmit={handleSubmit}
            onCancel={handleCancel}
            loading={loading}
            error={error ? 'An error occurred. Please try again.' : undefined}
          />
        )}
      </Box>
    </DashboardLayout>
  );
}

