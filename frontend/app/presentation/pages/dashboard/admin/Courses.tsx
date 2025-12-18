import React, { useState, useCallback } from 'react';
import {
  Box,
  Typography,
  IconButton,
  Card,
  CardContent,
  TablePagination,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import { useNavigate } from 'react-router';
import {
  useGetCoursesQuery,
  useDeleteCourseMutation,
} from '../../../../infrastructure/api/rtk/coursesApi';
import { useDashboardAuth } from '../../../hooks/useDashboardAuth';
import DashboardLayout from '../../../components/dashboard/Layout';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import Success from '../../../components/common/Success';
import Button from '../../../components/common/Button';
import { endpoints } from '../../../../infrastructure/api/endpoints';
import { ROUTES } from '../../../../shared/constants/routes';
import type { Course } from '../../../../domain/entities/Course';

const ITEMS_PER_PAGE = 10;

export default function AdminCourses() {
  const navigate = useNavigate();
  const { user: currentUser, logout } = useDashboardAuth();

  const [page, setPage] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [courseToDelete, setCourseToDelete] = useState<Course | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const {
    data,
    isLoading,
    error,
    refetch,
  } = useGetCoursesQuery({
    searchQuery: searchQuery || undefined,
    limit: ITEMS_PER_PAGE,
    offset: page * ITEMS_PER_PAGE,
    sortColumn: 'created_at',
    sortDirection: 'desc',
  });

  const [deleteCourse, { isLoading: deleting }] = useDeleteCourseMutation();

  const handleCreate = useCallback(() => {
    navigate(ROUTES.ADMIN_COURSE_DETAIL('new'));
  }, [navigate]);

  const handleEdit = useCallback((course: Course) => {
    navigate(ROUTES.ADMIN_COURSE_DETAIL(course.id));
  }, [navigate]);

  const handleDelete = useCallback((course: Course) => {
    setCourseToDelete(course);
    setDeleteDialogOpen(true);
  }, []);


  const handleConfirmDelete = useCallback(async () => {
    if (!courseToDelete) return;

    try {
      await deleteCourse(courseToDelete.id).unwrap();
      setSuccessMessage('Course deleted successfully');
      setDeleteDialogOpen(false);
      setCourseToDelete(null);
      refetch();
    } catch (err) {
      console.error('Delete failed:', err);
    }
  }, [courseToDelete, deleteCourse, refetch]);

  const handlePageChange = useCallback((_event: unknown, newPage: number) => {
    setPage(newPage);
  }, []);

  if (!currentUser || currentUser.role !== 'admin') {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h5" color="error">
          Access Denied. Admin role required.
        </Typography>
      </Box>
    );
  }

  const courses = data?.courses || [];
  const total = data?.total || 0;

  return (
    <DashboardLayout user={currentUser} onLogout={logout}>
      <Box>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1" sx={{ fontWeight: 600 }}>
            Course Management
          </Typography>
          <IconButton
            onClick={handleCreate}
            sx={{
              bgcolor: 'primary.main',
              color: 'primary.contrastText',
              '&:hover': { bgcolor: 'primary.dark' },
            }}
            aria-label="create course"
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
            <Box sx={{ mb: 2 }}>
              <input
                type="text"
                placeholder="Search courses..."
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setPage(0);
                }}
                style={{
                  width: '100%',
                  padding: '12px',
                  border: '1px solid #ccc',
                  borderRadius: '4px',
                  fontSize: '16px',
                }}
              />
            </Box>

            {error && (
              <Error
                message="Failed to load courses"
                onRetry={() => refetch()}
                fullWidth
              />
            )}

            {isLoading && <Loading variant="skeleton" fullWidth />}

            {!isLoading && !error && (
              <>
                <TableContainer component={Paper}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>No</TableCell>
                        <TableCell>Thumbnail</TableCell>
                        <TableCell>Name</TableCell>
                        <TableCell>Description</TableCell>
                        <TableCell>Categories</TableCell>
                        <TableCell>Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {courses.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={6} align="center">
                            <Typography variant="body2" color="text.secondary">
                              No courses found
                            </Typography>
                          </TableCell>
                        </TableRow>
                      ) : (
                        courses.map((course, index) => {
                          const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://asto-lms.local';
                          const thumbnailUrl = course.thumbnailUrl || 
                            (course.thumbnailId ? `${baseURL}${endpoints.buckets.download('course-thumbnails', course.thumbnailId)}` : null);
                          
                         
                          const rowNumber = page * ITEMS_PER_PAGE + index + 1;
                          
                       
                          if (import.meta.env.DEV) {
                            console.debug('[Courses] Course categories:', {
                              courseId: course.id,
                              courseName: course.name,
                              categories: course.categories,
                              categoriesLength: course.categories?.length,
                              categoriesType: typeof course.categories,
                              isArray: Array.isArray(course.categories),
                            });
                          }
                          
                          return (
                            <TableRow key={course.id} hover>
                              <TableCell>
                                <Typography variant="body2" color="text.secondary">
                                  {rowNumber}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                {thumbnailUrl ? (
                                  <img
                                    src={thumbnailUrl}
                                    alt={course.name}
                                    style={{
                                      width: '80px',
                                      height: '60px',
                                      objectFit: 'cover',
                                      borderRadius: '4px',
                                    }}
                                  />
                                ) : (
                                  <Box
                                    sx={{
                                      width: '80px',
                                      height: '60px',
                                      bgcolor: 'grey.200',
                                      display: 'flex',
                                      alignItems: 'center',
                                      justifyContent: 'center',
                                      borderRadius: '4px',
                                    }}
                                  >
                                    <Typography variant="caption" color="text.secondary">
                                      No image
                                    </Typography>
                                  </Box>
                                )}
                              </TableCell>
                              <TableCell>
                                <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                  {course.name}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                <Typography
                                  variant="body2"
                                  color="text.secondary"
                                  sx={{
                                    maxWidth: '300px',
                                    overflow: 'hidden',
                                    textOverflow: 'ellipsis',
                                    whiteSpace: 'nowrap',
                                  }}
                                >
                                  {course.description || 'No description'}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                                  {course.categories && Array.isArray(course.categories) && course.categories.length > 0 ? (
                                    course.categories.map((category) => {
                                  
                                      if (!category || !category.id || !category.name) {
                                        if (import.meta.env.DEV) {
                                          console.warn('[Courses] Invalid category object:', category);
                                        }
                                        return null;
                                      }
                                      return (
                                        <Chip
                                          key={category.id}
                                          label={category.name}
                                          size="small"
                                          color="primary"
                                          variant="outlined"
                                        />
                                      );
                                    }).filter(Boolean)
                                  ) : (
                                    <Typography variant="caption" color="text.secondary">
                                      No categories
                                    </Typography>
                                  )}
                                </Box>
                              </TableCell>
                              <TableCell>
                                <Box sx={{ display: 'flex', gap: 1 }}>
                                  <IconButton
                                    size="small"
                                    onClick={() => handleEdit(course)}
                                    color="primary"
                                  >
                                    <EditIcon />
                                  </IconButton>
                                  <IconButton
                                    size="small"
                                    onClick={() => handleDelete(course)}
                                    color="error"
                                  >
                                    <DeleteIcon />
                                  </IconButton>
                                </Box>
                              </TableCell>
                            </TableRow>
                          );
                        })
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>
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
          </CardContent>
        </Card>

        <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
          <DialogTitle>Delete Course</DialogTitle>
          <DialogContent>
            <DialogContentText>
              Are you sure you want to delete course &quot;{courseToDelete?.name}&quot;? This action
              cannot be undone.
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDeleteDialogOpen(false)} variant="outline" disabled={deleting}>
              Cancel
            </Button>
            <Button
              onClick={handleConfirmDelete}
              variant="fill"
              color="secondary"
              disabled={deleting}
              sx={{ bgcolor: 'error.main', '&:hover': { bgcolor: 'error.dark' } }}
            >
              {deleting ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </DashboardLayout>
  );
}

