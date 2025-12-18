import React, { useCallback } from 'react';
import { useNavigate } from 'react-router';
import { useDispatch } from 'react-redux';
import {
  Box,
  Typography,
  Chip,
  Card,
  CardContent,
  CardActions,
  Grid,
} from '@mui/material';
import { useAuth } from '../../../hooks/useAuth';
import { useGetEnrollmentsQuery } from '../../../../infrastructure/api/rtk/enrollmentsApi';
import { clearPublicAuth, clearDashboardAuth } from '../../../../infrastructure/store/authSlice';
import DashboardLayout from '../../../components/dashboard/Layout';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import Button from '../../../components/common/Button';
import { ROUTES } from '../../../../shared/constants/routes';

export default function StudentDashboard() {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { user, loading: authLoading } = useAuth();

  const studentId = user?.id;

  const {
    data,
    isLoading,
    error,
    refetch,
  } = useGetEnrollmentsQuery(
    studentId
      ? {
          studentId,
          limit: 50,
          offset: 0,
          sortColumn: 'created_at',
          sortDirection: 'desc',
        }
      : undefined,
    {
      skip: !studentId,
    }
  );

  const enrollments = data?.enrollments || [];

  const handleGoToClass = (courseOfferingId: string) => {
    navigate(ROUTES.STUDENT_CLASS(courseOfferingId));
  };

  const handleLogout = useCallback(() => {
    dispatch(clearPublicAuth());
    dispatch(clearDashboardAuth());
    navigate(ROUTES.HOME);
  }, [dispatch, navigate]);

  if (authLoading || (!user && authLoading)) {
    return <Loading variant="fullscreen" />;
  }

  if (!user || user.role !== 'student') {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h5" color="error">
          Access Denied. Student role required.
        </Typography>
      </Box>
    );
  }

  return (
    <DashboardLayout user={user} onLogout={handleLogout}>
      <Box sx={{ px: { xs: 2, md: 4 }, py: 4 }}>
        <Typography
          variant="h4"
          sx={{
            mb: 3,
            fontWeight: 700,
          }}
        >
          My Classes
        </Typography>

        {error && (
          <Error
            message="Failed to load your enrollments."
            onRetry={() => refetch()}
            fullWidth
          />
        )}

        {isLoading && <Loading variant="skeleton" fullWidth />}

        {!isLoading && !error && enrollments.length === 0 && (
          <Box sx={{ py: 8, textAlign: 'center', color: 'text.secondary' }}>
            <Typography variant="h6" sx={{ mb: 1 }}>
              You’re not enrolled in any classes yet.
            </Typography>
            <Typography variant="body2">
              Browse available courses on the public site and enroll to see your classes here.
            </Typography>
          </Box>
        )}

        {!isLoading && !error && enrollments.length > 0 && (
          <Grid container spacing={3}>
            {enrollments.map((enrollment) => {
              const canEnter =
                enrollment.status === 'approved' || enrollment.status === 'completed';

              return (
                <Grid item xs={12} md={6} key={enrollment.id}>
                  <Card
                    sx={{
                      height: '100%',
                      display: 'flex',
                      flexDirection: 'column',
                      borderRadius: 3,
                      boxShadow: 2,
                    }}
                  >
                    <CardContent
                      sx={{
                        display: 'flex',
                        flexDirection: 'column',
                        gap: 1.5,
                      }}
                    >
                      <Typography variant="subtitle2" color="text.secondary">
                        {enrollment.courseName}
                      </Typography>
                      <Typography variant="h6" sx={{ fontWeight: 600 }}>
                        {enrollment.courseOfferingName}
                      </Typography>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
                        <Chip
                          label={enrollment.status}
                          size="small"
                          color={
                            enrollment.status === 'approved' || enrollment.status === 'completed'
                              ? 'success'
                              : enrollment.status === 'pending'
                              ? 'warning'
                              : 'default'
                          }
                          sx={{ textTransform: 'capitalize' }}
                        />
                        <Typography variant="caption" color="text.secondary">
                          Enrolled on {new Date(enrollment.createdAt).toLocaleDateString()}
                        </Typography>
                      </Box>
                    </CardContent>
                    <CardActions sx={{ px: 2.5, pb: 2.5, pt: 0, mt: 'auto' }}>
                      {canEnter ? (
                        <Button
                          variant="fill"
                          size="md"
                          onClick={() => handleGoToClass(enrollment.courseOfferingId)}
                        >
                          Go to class
                        </Button>
                      ) : (
                        <Typography variant="body2" color="text.secondary">
                          {enrollment.status === 'pending'
                            ? 'Waiting for approval from the academy.'
                            : 'This enrollment is not active.'}
                        </Typography>
                      )}
                    </CardActions>
                  </Card>
                </Grid>
              );
            })}
          </Grid>
        )}
      </Box>
    </DashboardLayout>
  );
}


