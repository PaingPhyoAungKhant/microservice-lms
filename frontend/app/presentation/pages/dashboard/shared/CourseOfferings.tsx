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
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import { useNavigate, useSearchParams } from 'react-router';
import {
  useGetCourseOfferingsQuery,
} from '../../../../infrastructure/api/rtk/courseOfferingsApi';
import { useDashboardAuth } from '../../../hooks/useDashboardAuth';
import DashboardLayout from '../../../components/dashboard/Layout';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import type { CourseOffering } from '../../../../domain/entities/CourseOffering';

const ITEMS_PER_PAGE = 10;

interface CourseOfferingsProps {
  routePrefix: string;
  allowedRoles: ('admin' | 'instructor')[];
}

export default function CourseOfferings({ routePrefix, allowedRoles }: CourseOfferingsProps) {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const courseId = searchParams.get('courseId');
  const { user: currentUser, logout } = useDashboardAuth();

  const [page, setPage] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');

  const {
    data,
    isLoading,
    error,
    refetch,
  } = useGetCourseOfferingsQuery({
    search: searchQuery || undefined,
    courseId: courseId || undefined,
    limit: ITEMS_PER_PAGE,
    offset: page * ITEMS_PER_PAGE,
    sortColumn: 'created_at',
    sortDirection: 'desc',
  });

  const handleCreate = useCallback(() => {
    if (courseId) {
      navigate(`${routePrefix}/course-offerings/new?courseId=${courseId}`);
    } else {
      navigate(`${routePrefix}/course-offerings/new`);
    }
  }, [navigate, courseId, routePrefix]);

  const handleEdit = useCallback((offering: CourseOffering) => {
    navigate(`${routePrefix}/course-offerings/${offering.id}`);
  }, [navigate, routePrefix]);

  const handlePageChange = useCallback((_event: unknown, newPage: number) => {
    setPage(newPage);
  }, []);

  if (!currentUser || !allowedRoles.includes(currentUser.role as 'admin' | 'instructor')) {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h5" color="error">
          Access Denied. {allowedRoles.map(r => r.charAt(0).toUpperCase() + r.slice(1)).join(' or ')} role required.
        </Typography>
      </Box>
    );
  }

  const offerings = data?.courseOfferings || [];
  const total = data?.total || 0;

  return (
    <DashboardLayout user={currentUser} onLogout={logout}>
      <Box>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1" sx={{ fontWeight: 600 }}>
            Course Offerings Management
          </Typography>
          <IconButton
            onClick={handleCreate}
            sx={{
              bgcolor: 'primary.main',
              color: 'primary.contrastText',
              '&:hover': { bgcolor: 'primary.dark' },
            }}
            aria-label="create course offering"
          >
            <AddIcon />
          </IconButton>
        </Box>

        <Card>
          <CardContent>
            <Box sx={{ mb: 2 }}>
              <input
                type="text"
                placeholder="Search course offerings..."
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
                message="Failed to load course offerings"
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
                        <TableCell>Name</TableCell>
                        <TableCell>Course</TableCell>
                        <TableCell>Type</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Enrollment Cost</TableCell>
                        <TableCell>Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {offerings.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={7} align="center">
                            <Typography variant="body2" color="text.secondary">
                              No course offerings found
                            </Typography>
                          </TableCell>
                        </TableRow>
                      ) : (
                        offerings.map((offering, index) => {
                          const rowNumber = page * ITEMS_PER_PAGE + index + 1;
                          
                          return (
                            <TableRow key={offering.id} hover>
                              <TableCell>
                                <Typography variant="body2" color="text.secondary">
                                  {rowNumber}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                <Typography variant="body1" sx={{ fontWeight: 500 }}>
                                  {offering.name}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                <Typography variant="body2" color="text.secondary">
                                  {offering.courseName || 'N/A'}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                <Chip
                                  label={offering.offeringType}
                                  size="small"
                                  color={offering.offeringType === 'online' ? 'primary' : 'secondary'}
                                  variant="outlined"
                                />
                              </TableCell>
                              <TableCell>
                                <Chip
                                  label={offering.status}
                                  size="small"
                                  color={
                                    offering.status === 'active' ? 'success' :
                                    offering.status === 'pending' ? 'warning' :
                                    offering.status === 'completed' ? 'default' : 'info'
                                  }
                                  variant="outlined"
                                />
                              </TableCell>
                              <TableCell>
                                <Typography variant="body2">
                                  ${offering.enrollmentCost.toFixed(2)}
                                </Typography>
                              </TableCell>
                              <TableCell>
                                <Box sx={{ display: 'flex', gap: 1 }}>
                                  <IconButton
                                    size="small"
                                    onClick={() => handleEdit(offering)}
                                    color="primary"
                                  >
                                    <EditIcon />
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
      </Box>
    </DashboardLayout>
  );
}

