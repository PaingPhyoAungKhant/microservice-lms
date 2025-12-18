import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Checkbox,
  FormControlLabel,
  List,
  ListItem,
  Paper,
  TextField,
} from '@mui/material';
import Button from '../../../components/common/Button';
import type { User } from '../../../../domain/entities/User';
import type { CourseOffering } from '../../../../domain/entities/CourseOffering';

interface EnrollmentFormProps {
  open: boolean;
  courseOffering: CourseOffering | null;
  students: User[];
  enrolledStudentIds: string[];
  onClose: () => void;
  onSubmit: (studentIds: string[]) => Promise<void>;
  loading?: boolean;
}

export default function EnrollmentForm({
  open,
  courseOffering,
  students,
  enrolledStudentIds,
  onClose,
  onSubmit,
  loading = false,
}: EnrollmentFormProps) {
  const [selectedStudentIds, setSelectedStudentIds] = useState<string[]>([]);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    if (!open) {
      setSelectedStudentIds([]);
      setSearchQuery('');
    }
  }, [open]);

  const availableStudents = students.filter(
    (student) =>
      !enrolledStudentIds.includes(student.id) &&
      student.role === 'student' &&
      (student.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
        student.email.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  const handleToggleStudent = (studentId: string) => {
    setSelectedStudentIds((prev) =>
      prev.includes(studentId)
        ? prev.filter((id) => id !== studentId)
        : [...prev, studentId]
    );
  };

  const handleSelectAll = () => {
    if (selectedStudentIds.length === availableStudents.length) {
      setSelectedStudentIds([]);
    } else {
      setSelectedStudentIds(availableStudents.map((s) => s.id));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedStudentIds.length === 0) {
      return;
    }
    try {
      await onSubmit(selectedStudentIds);
      onClose();
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        Enroll Students
        {courseOffering && (
          <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
            {courseOffering.name}
          </Typography>
        )}
      </DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <TextField
            fullWidth
            label="Search students"
            variant="outlined"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            sx={{ mb: 2 }}
            placeholder="Search by username or email..."
          />
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
            <Typography variant="body2" color="text.secondary">
              {availableStudents.length} available student(s)
            </Typography>
            {availableStudents.length > 0 && (
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectedStudentIds.length === availableStudents.length}
                    indeterminate={
                      selectedStudentIds.length > 0 &&
                      selectedStudentIds.length < availableStudents.length
                    }
                    onChange={handleSelectAll}
                  />
                }
                label="Select All"
              />
            )}
          </Box>
          <Paper
            variant="outlined"
            sx={{
              maxHeight: 400,
              overflow: 'auto',
              minHeight: 200,
            }}
          >
            {availableStudents.length === 0 ? (
              <Box sx={{ textAlign: 'center', py: 4 }}>
                <Typography color="text.secondary">
                  {searchQuery
                    ? 'No students found matching your search'
                    : 'No available students to enroll'}
                </Typography>
              </Box>
            ) : (
              <List>
                {availableStudents.map((student) => (
                  <ListItem key={student.id} disablePadding>
                    <FormControlLabel
                      control={
                        <Checkbox
                          checked={selectedStudentIds.includes(student.id)}
                          onChange={() => handleToggleStudent(student.id)}
                        />
                      }
                      label={
                        <Box>
                          <Typography variant="body1">{student.username}</Typography>
                          <Typography variant="caption" color="text.secondary">
                            {student.email}
                          </Typography>
                        </Box>
                      }
                      sx={{ width: '100%', px: 2 }}
                    />
                  </ListItem>
                ))}
              </List>
            )}
          </Paper>
          {selectedStudentIds.length > 0 && (
            <Typography variant="body2" color="primary" sx={{ mt: 2 }}>
              {selectedStudentIds.length} student(s) selected
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} variant="outline" disabled={loading}>
            Cancel
          </Button>
          <Button
            type="submit"
            variant="fill"
            color="primary"
            disabled={loading || selectedStudentIds.length === 0}
          >
            {loading ? 'Enrolling...' : `Enroll ${selectedStudentIds.length} Student(s)`}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}

