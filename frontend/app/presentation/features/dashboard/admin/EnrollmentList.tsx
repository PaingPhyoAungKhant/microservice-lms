import React from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Chip,
  Box,
  Typography,
  Select,
  MenuItem,
  FormControl,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import type { Enrollment, EnrollmentStatus } from '../../../../domain/entities/Enrollment';

interface EnrollmentListProps {
  enrollments: Enrollment[];
  loading: boolean;
  onDelete: (enrollment: Enrollment) => void;
  onStatusUpdate: (enrollmentId: string, status: EnrollmentStatus) => void;
  updatingStatus?: string | null;
  page: number;
  itemsPerPage: number;
}

export default function EnrollmentList({
  enrollments,
  loading,
  onDelete,
  onStatusUpdate,
  updatingStatus,
  page,
  itemsPerPage,
}: EnrollmentListProps) {
  const getStatusColor = (
    status: string
  ): 'success' | 'warning' | 'error' | 'default' | 'info' => {
    switch (status) {
      case 'approved':
        return 'success';
      case 'pending':
        return 'warning';
      case 'rejected':
        return 'error';
      case 'completed':
        return 'info';
      default:
        return 'default';
    }
  };

  if (loading) {
    return (
      <Box sx={{ textAlign: 'center', py: 4 }}>
        <Typography color="text.secondary">Loading enrollments...</Typography>
      </Box>
    );
  }

  if (enrollments.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 4 }}>
        <Typography color="text.secondary">No enrollments found</Typography>
      </Box>
    );
  }

  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>No.</TableCell>
            <TableCell>Student</TableCell>
            <TableCell>Course</TableCell>
            <TableCell>Course Offering</TableCell>
            <TableCell>Status</TableCell>
            <TableCell align="right">Actions</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {enrollments.map((enrollment, index) => (
            <TableRow key={enrollment.id} hover>
              <TableCell>{page * itemsPerPage + index + 1}</TableCell>
              <TableCell>{enrollment.studentUsername}</TableCell>
              <TableCell>{enrollment.courseName}</TableCell>
              <TableCell>{enrollment.courseOfferingName}</TableCell>
              <TableCell>
                <FormControl size="small" sx={{ minWidth: 120 }}>
                  <Select
                    value={enrollment.status}
                    onChange={(e) =>
                      onStatusUpdate(enrollment.id, e.target.value as EnrollmentStatus)
                    }
                    disabled={updatingStatus === enrollment.id}
                    sx={{
                      '& .MuiSelect-select': {
                        py: 0.5,
                      },
                    }}
                  >
                    <MenuItem value="pending">Pending</MenuItem>
                    <MenuItem value="approved">Approved</MenuItem>
                    <MenuItem value="rejected">Rejected</MenuItem>
                    <MenuItem value="completed">Completed</MenuItem>
                  </Select>
                </FormControl>
              </TableCell>
              <TableCell align="right">
                <IconButton
                  size="small"
                  onClick={() => onDelete(enrollment)}
                  color="error"
                  aria-label="delete enrollment"
                >
                  <DeleteIcon />
                </IconButton>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

