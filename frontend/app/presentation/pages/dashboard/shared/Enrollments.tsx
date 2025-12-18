import React, { useState, useCallback } from 'react';
import {
  Box,
  Typography,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
  TablePagination,
  Card,
  CardContent,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { useSearchParams } from 'react-router';
import {
  useGetEnrollmentsQuery,
  useCreateEnrollmentMutation,
  useUpdateEnrollmentStatusMutation,
  useDeleteEnrollmentMutation,
} from '../../../../infrastructure/api/rtk/enrollmentsApi';
import { useGetCourseOfferingsQuery } from '../../../../infrastructure/api/rtk/courseOfferingsApi';
import { useUsers } from '../../../hooks/useUsers';
import { useDashboardAuth } from '../../../hooks/useDashboardAuth';
import DashboardLayout from '../../../components/dashboard/Layout';
import EnrollmentList from '../../../features/dashboard/admin/EnrollmentList';
import EnrollmentForm from '../../../features/dashboard/admin/EnrollmentForm';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import Success from '../../../components/common/Success';
import Button from '../../../components/common/Button';
import type { Enrollment, EnrollmentStatus } from '../../../../domain/entities/Enrollment';
import type { CourseOffering } from '../../../../domain/entities/CourseOffering';

const ITEMS_PER_PAGE = 10;

interface EnrollmentsProps {
  routePrefix: string;
  allowedRoles: ('admin' | 'instructor')[];
}

export default function Enrollments({ routePrefix, allowedRoles }: EnrollmentsProps) {
  const { user: currentUser, logout } = useDashboardAuth();
  const [searchParams, setSearchParams] = useSearchParams();
  const courseOfferingIdParam = searchParams.get('courseOfferingId');

  const [page, setPage] = useState(0);
  const [selectedCourseOfferingId, setSelectedCourseOfferingId] = useState<string>(
    courseOfferingIdParam || ''
  );
  const [formOpen, setFormOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [enrollmentToDelete, setEnrollmentToDelete] = useState<Enrollment | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [updatingStatusId, setUpdatingStatusId] = useState<string | null>(null);

  const { data: courseOfferingsData } = useGetCourseOfferingsQuery({});

  const {
    data: enrollmentsData,
    isLoading: enrollmentsLoading,
    error: enrollmentsError,
    refetch: refetchEnrollments,
  } = useGetEnrollmentsQuery({
    courseOfferingId: selectedCourseOfferingId || undefined,
    limit: ITEMS_PER_PAGE,
    offset: page * ITEMS_PER_PAGE,
    sortColumn: 'created_at',
    sortDirection: 'desc',
  });

  const { users: students } = useUsers({
    role: 'student',
    status: 'active',
    limit: 100,
  });

  const [createEnrollment, { isLoading: creating }] = useCreateEnrollmentMutation();
  const [updateEnrollmentStatus, { isLoading: updatingStatus }] =
    useUpdateEnrollmentStatusMutation();
  const [deleteEnrollment, { isLoading: deleting }] = useDeleteEnrollmentMutation();

  const enrollments = enrollmentsData?.enrollments || [];
  const total = enrollmentsData?.total || 0;
  const courseOfferings = courseOfferingsData?.courseOfferings || [];
  const selectedCourseOffering = courseOfferings.find(
    (co) => co.id === selectedCourseOfferingId
  );

 
  const enrolledStudentIds = enrollments.map((e) => e.studentId);

  const handleCourseOfferingChange = (courseOfferingId: string) => {
    setSelectedCourseOfferingId(courseOfferingId);
    setPage(0);
    if (courseOfferingId) {
      setSearchParams({ courseOfferingId });
    } else {
      setSearchParams({});
    }
  };

  const handleCreate = () => {
    if (!selectedCourseOfferingId) {
      setSuccessMessage('Please select a course offering first');
      return;
    }
    setFormOpen(true);
  };

  const handleDelete = (enrollment: Enrollment) => {
    setEnrollmentToDelete(enrollment);
    setDeleteDialogOpen(true);
  };

  const handleFormSubmit = async (studentIds: string[]) => {
    if (!selectedCourseOffering || !selectedCourseOfferingId) {
      throw new Error('Course offering not selected');
    }

    try {
      const courseId = selectedCourseOffering.courseId;
      const courseName = selectedCourseOffering.courseName || 'Unknown Course';

      const enrollPromises = studentIds.map(async (studentId) => {
        const student = students.find((s) => s.id === studentId);
        if (!student) {
          throw new Error(`Student with ID ${studentId} not found`);
        }

        return createEnrollment({
          studentId: student.id,
          studentUsername: student.username,
          courseId: courseId,
          courseName: courseName,
          courseOfferingId: selectedCourseOfferingId,
          courseOfferingName: selectedCourseOffering.name,
        }).unwrap();
      });

      await Promise.all(enrollPromises);
      setSuccessMessage(`${studentIds.length} student(s) enrolled successfully`);
      setFormOpen(false);
      refetchEnrollments();
    } catch (err) {
      console.error('Failed to enroll students:', err);
      throw err;
    }
  };

  const handleStatusUpdate = async (enrollmentId: string, status: EnrollmentStatus) => {
    setUpdatingStatusId(enrollmentId);
    try {
      await updateEnrollmentStatus({ id: enrollmentId, status }).unwrap();
      setSuccessMessage('Enrollment status updated successfully');
      refetchEnrollments();
    } catch (err) {
      console.error('Failed to update enrollment status:', err);
      setSuccessMessage('Failed to update enrollment status');
    } finally {
      setUpdatingStatusId(null);
    }
  };

  const handleConfirmDelete = async () => {
    if (!enrollmentToDelete) return;

    try {
      await deleteEnrollment(enrollmentToDelete.id).unwrap();
      setSuccessMessage('Enrollment removed successfully');
      setDeleteDialogOpen(false);
      setEnrollmentToDelete(null);
      refetchEnrollments();
    } catch (err) {
      console.error('Failed to delete enrollment:', err);
    }
  };

  const handlePageChange = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  if (!currentUser || !allowedRoles.includes(currentUser.role as 'admin' | 'instructor')) {
    return (
      <DashboardLayout user={currentUser} onLogout={logout}>
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h5" color="error">
            Access Denied. {allowedRoles.map(r => r.charAt(0).toUpperCase() + r.slice(1)).join(' or ')} role required.
          </Typography>
        </Box>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout user={currentUser} onLogout={logout}>
      <Box>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1" sx={{ fontWeight: 600 }}>
            Enrollment Management
          </Typography>
          <IconButton
            onClick={handleCreate}
            disabled={!selectedCourseOfferingId}
            sx={{
              bgcolor: 'primary.main',
              color: 'primary.contrastText',
              '&:hover': { bgcolor: 'primary.dark' },
              '&:disabled': { bgcolor: 'action.disabledBackground' },
            }}
            aria-label="enroll students"
          >
            <AddIcon />
          </IconButton>
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

        <Card>
          <CardContent>
            <Box sx={{ mb: 3 }}>
              <FormControl fullWidth>
                <InputLabel>Course Offering</InputLabel>
                <Select
                  value={selectedCourseOfferingId}
                  label="Course Offering"
                  onChange={(e) => handleCourseOfferingChange(e.target.value)}
                >
                  <MenuItem value="">
                    <em>Select a course offering</em>
                  </MenuItem>
                  {courseOfferings.map((offering) => (
                    <MenuItem key={offering.id} value={offering.id}>
                      {offering.name} - {offering.courseName || 'N/A'}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Box>

            {!selectedCourseOfferingId && (
              <Box sx={{ textAlign: 'center', py: 4 }}>
                <Typography variant="body2" color="text.secondary">
                  Please select a course offering to view enrollments
                </Typography>
              </Box>
            )}

            {selectedCourseOfferingId && (
              <>
                {enrollmentsError && (
                  <Error
                    message="Failed to load enrollments"
                    onRetry={() => refetchEnrollments()}
                    fullWidth
                  />
                )}

                {enrollmentsLoading && <Loading variant="skeleton" fullWidth />}

                {!enrollmentsLoading && !enrollmentsError && (
                  <>
                    <EnrollmentList
                      enrollments={enrollments}
                      loading={enrollmentsLoading}
                      onDelete={handleDelete}
                      onStatusUpdate={handleStatusUpdate}
                      updatingStatus={updatingStatusId}
                      page={page}
                      itemsPerPage={ITEMS_PER_PAGE}
                    />
                    <TablePagination
                      component="div"
                      count={total}
                      page={page}
                      onPageChange={handlePageChange}
                      rowsPerPage={ITEMS_PER_PAGE}
                      rowsPerPageOptions={[]}
                    />
                  </>
                )}
              </>
            )}
          </CardContent>
        </Card>

        <EnrollmentForm
          open={formOpen}
          courseOffering={selectedCourseOffering || null}
          students={students}
          enrolledStudentIds={enrolledStudentIds}
          onClose={() => setFormOpen(false)}
          onSubmit={handleFormSubmit}
          loading={creating}
        />

        <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
          <DialogTitle>Remove Enrollment</DialogTitle>
          <DialogContent>
            <DialogContentText>
              Are you sure you want to remove &quot;{enrollmentToDelete?.studentUsername}&quot; from
              &quot;{enrollmentToDelete?.courseOfferingName}&quot;? This action cannot be undone.
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDeleteDialogOpen(false)} variant="outline" disabled={deleting}>
              Cancel
            </Button>
            <Button
              onClick={handleConfirmDelete}
              variant="fill"
              color="error"
              disabled={deleting}
            >
              {deleting ? 'Removing...' : 'Remove'}
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </DashboardLayout>
  );
}

