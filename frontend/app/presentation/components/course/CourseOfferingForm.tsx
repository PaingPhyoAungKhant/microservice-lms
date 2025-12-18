import React, { useState, useEffect } from 'react';
import {
  Box,
  TextField,
  Typography,
  Paper,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  FormHelperText,
} from '@mui/material';
import { useGetCoursesQuery } from '../../../infrastructure/api/rtk/coursesApi';
import Button from '../common/Button';
import type { CourseOffering } from '../../../domain/entities/CourseOffering';

interface CourseOfferingFormProps {
  courseOffering?: CourseOffering | null;
  courseId?: string | null;
  onSubmit: (data: {
    courseId: string;
    name: string;
    description: string;
    offeringType: 'online' | 'oncampus';
    duration?: string | null;
    classTime?: string | null;
    enrollmentCost: number;
    status?: 'pending' | 'active' | 'ongoing' | 'completed';
  }) => void;
  onCancel?: () => void;
  loading?: boolean;
  error?: string;
}

export default function CourseOfferingForm({
  courseOffering,
  courseId: initialCourseId,
  onSubmit,
  onCancel,
  loading = false,
  error,
}: CourseOfferingFormProps) {
  const [name, setName] = useState(courseOffering?.name || '');
  const [description, setDescription] = useState(courseOffering?.description || '');
  const [offeringType, setOfferingType] = useState<'online' | 'oncampus'>(
    (courseOffering?.offeringType as 'online' | 'oncampus') || 'online'
  );
  const [duration, setDuration] = useState(courseOffering?.duration || '');
  const [classTime, setClassTime] = useState(courseOffering?.classTime || '');
  const [enrollmentCost, setEnrollmentCost] = useState(
    courseOffering?.enrollmentCost.toString() || '0'
  );
  const [status, setStatus] = useState<'pending' | 'active' | 'ongoing' | 'completed'>(
    (courseOffering?.status as 'pending' | 'active' | 'ongoing' | 'completed') || 'pending'
  );
  const [selectedCourseId, setSelectedCourseId] = useState<string>(initialCourseId || courseOffering?.courseId || '');
  const [nameError, setNameError] = useState('');
  const [costError, setCostError] = useState('');
  const [courseError, setCourseError] = useState('');

  const { data: coursesData, isLoading: coursesLoading } = useGetCoursesQuery({
    limit: 1000,
    sortColumn: 'name',
    sortDirection: 'asc',
  });

  const courses = coursesData?.courses || [];

  useEffect(() => {
    if (courseOffering) {
      setName(courseOffering.name);
      setDescription(courseOffering.description);
      setOfferingType(courseOffering.offeringType);
      setDuration(courseOffering.duration || '');
      setClassTime(courseOffering.classTime || '');
      setEnrollmentCost(courseOffering.enrollmentCost.toString());
      setStatus(courseOffering.status as 'pending' | 'active' | 'ongoing' | 'completed');
      setSelectedCourseId(courseOffering.courseId);
    } else if (initialCourseId) {
      setSelectedCourseId(initialCourseId);
    }
  }, [courseOffering, initialCourseId]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

  
    let hasError = false;
    if (!name.trim()) {
      setNameError('Course offering name is required');
      hasError = true;
    }
    if (!selectedCourseId) {
      setCourseError('Course is required');
      hasError = true;
    }
    const cost = parseFloat(enrollmentCost);
    if (isNaN(cost) || cost < 0) {
      setCostError('Enrollment cost must be a valid number >= 0');
      hasError = true;
    }
    if (hasError) {
      return;
    }

    onSubmit({
      courseId: selectedCourseId,
      name: name.trim(),
      description: description.trim(),
      offeringType,
      duration: duration.trim() || null,
      classTime: classTime.trim() || null,
      enrollmentCost: cost,
      ...(courseOffering && { status }),
    });
  };

  return (
    <form onSubmit={handleSubmit}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          {courseOffering ? 'Edit Course Offering' : 'Create Course Offering'}
        </Typography>

        {error && (
          <Box sx={{ mb: 2, p: 2, bgcolor: 'error.light', borderRadius: 1 }}>
            <Typography variant="body2" color="error">
              {error}
            </Typography>
          </Box>
        )}

        <FormControl fullWidth sx={{ mb: 2 }} required error={!!courseError}>
          <InputLabel>Course</InputLabel>
          <Select
            value={selectedCourseId}
            onChange={(e) => {
              setSelectedCourseId(e.target.value);
              setCourseError('');
            }}
            label="Course"
            disabled={!!initialCourseId || !!courseOffering}
          >
            {coursesLoading ? (
              <MenuItem disabled>Loading courses...</MenuItem>
            ) : courses.length === 0 ? (
              <MenuItem disabled>No courses available</MenuItem>
            ) : (
              courses.map((course) => (
                <MenuItem key={course.id} value={course.id}>
                  {course.name}
                </MenuItem>
              ))
            )}
          </Select>
          {courseError && <FormHelperText>{courseError}</FormHelperText>}
        </FormControl>

        <TextField
          fullWidth
          label="Course Offering Name"
          value={name}
          onChange={(e) => {
            setName(e.target.value);
            setNameError('');
          }}
          error={!!nameError}
          helperText={nameError}
          required
          sx={{ mb: 2 }}
        />

        <TextField
          fullWidth
          label="Description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          multiline
          rows={4}
          sx={{ mb: 2 }}
        />

        <FormControl fullWidth sx={{ mb: 2 }} required>
          <InputLabel>Offering Type</InputLabel>
          <Select
            value={offeringType}
            onChange={(e) => setOfferingType(e.target.value as 'online' | 'oncampus')}
            label="Offering Type"
          >
            <MenuItem value="online">Online</MenuItem>
            <MenuItem value="oncampus">On Campus</MenuItem>
          </Select>
        </FormControl>

        <TextField
          fullWidth
          label="Duration"
          value={duration}
          onChange={(e) => setDuration(e.target.value)}
          placeholder="e.g., 12 weeks, 3 months"
          sx={{ mb: 2 }}
        />

        <TextField
          fullWidth
          label="Class Time"
          value={classTime}
          onChange={(e) => setClassTime(e.target.value)}
          placeholder="e.g., Monday, Wednesday, Friday 10:00 AM - 12:00 PM"
          sx={{ mb: 2 }}
        />

        <TextField
          fullWidth
          label="Enrollment Cost"
          type="number"
          value={enrollmentCost}
          onChange={(e) => {
            setEnrollmentCost(e.target.value);
            setCostError('');
          }}
          error={!!costError}
          helperText={costError}
          required
          inputProps={{ min: 0, step: 0.01 }}
          sx={{ mb: 2 }}
        />

        {courseOffering && (
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Status</InputLabel>
            <Select
              value={status}
              onChange={(e) => setStatus(e.target.value as 'pending' | 'active' | 'ongoing' | 'completed')}
              label="Status"
            >
              <MenuItem value="pending">Pending</MenuItem>
              <MenuItem value="active">Active</MenuItem>
              <MenuItem value="ongoing">Ongoing</MenuItem>
              <MenuItem value="completed">Completed</MenuItem>
            </Select>
          </FormControl>
        )}

        <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 3 }}>
          {onCancel && (
            <Button variant="outline" onClick={onCancel} disabled={loading}>
              Cancel
            </Button>
          )}
          <Button variant="fill" type="submit" disabled={loading || !selectedCourseId}>
            {loading ? 'Saving...' : courseOffering ? 'Update Course Offering' : 'Create Course Offering'}
          </Button>
        </Box>
      </Paper>
    </form>
  );
}

