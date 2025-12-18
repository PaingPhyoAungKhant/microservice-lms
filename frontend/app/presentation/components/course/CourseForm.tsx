import React, { useState, useEffect } from 'react';
import {
  Box,
  TextField,
  Typography,
  Paper,
} from '@mui/material';
import ImageIcon from '@mui/icons-material/Image';
import { useGetFileQuery } from '../../../infrastructure/api/rtk/filesApi';
import { getFileDownloadUrl } from '../../../infrastructure/api/utils';
import Button from '../common/Button';
import FileSelector from '../file/FileSelector';
import CategorySelector from './CategorySelector';
import type { Course } from '../../../domain/entities/Course';
import type { Category } from '../../../domain/entities/Category';

interface CourseFormProps {
  course?: Course | null;
  onSubmit: (data: {
    name: string;
    description: string;
    thumbnailId?: string | null;
    categoryIds: string[];
  }) => void;
  onCancel?: () => void;
  loading?: boolean;
  error?: string;
}

export default function CourseForm({
  course,
  onSubmit,
  onCancel,
  loading = false,
  error,
}: CourseFormProps) {
  const [name, setName] = useState(course?.name || '');
  const [description, setDescription] = useState(course?.description || '');
  const [thumbnailId, setThumbnailId] = useState<string | null>(course?.thumbnailId || null);
  const [selectedCategories, setSelectedCategories] = useState<Category[]>(
    course?.categories || []
  );
  const [showFileSelector, setShowFileSelector] = useState(false);
  const [nameError, setNameError] = useState('');

  const { data: thumbnailFile } = useGetFileQuery(thumbnailId || '', {
    skip: !thumbnailId,
  });

  useEffect(() => {
    if (course) {
      setName(course.name);
      setDescription(course.description);
      setThumbnailId(course.thumbnailId || null);
      setSelectedCategories(course.categories || []);
    }
  }, [course]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      setNameError('Course name is required');
      return;
    }

    onSubmit({
      name: name.trim(),
      description: description.trim(),
      thumbnailId: thumbnailId || null,
      categoryIds: selectedCategories.map((cat) => cat.id),
    });
  };

  const handleThumbnailSelect = (fileId: string) => {
    setThumbnailId(fileId);
  };

  const getThumbnailUrl = () => {
    if (course?.thumbnailUrl) {
      return course.thumbnailUrl;
    }
    if (thumbnailId) {
      return getFileDownloadUrl('course-thumbnails', thumbnailId);
    }
    return null;
  };

  const thumbnailUrl = getThumbnailUrl();

  return (
    <form onSubmit={handleSubmit}>
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          {course ? 'Edit Course' : 'Create Course'}
        </Typography>

        {error && (
          <Box sx={{ mb: 2, p: 2, bgcolor: 'error.light', borderRadius: 1 }}>
            <Typography variant="body2" color="error">
              {error}
            </Typography>
          </Box>
        )}

        <TextField
          fullWidth
          label="Course Name"
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

        <Box sx={{ mb: 2 }}>
          <Typography variant="subtitle2" gutterBottom sx={{ fontWeight: 600, mb: 1 }}>
            Thumbnail
          </Typography>
          {thumbnailUrl ? (
            <Box sx={{ mb: 2 }}>
              <Box
                sx={{
                  position: 'relative',
                  display: 'inline-block',
                  border: '1px solid',
                  borderColor: 'divider',
                  borderRadius: 1,
                  overflow: 'hidden',
                }}
              >
                <img
                  src={thumbnailUrl}
                  alt="Course thumbnail"
                  style={{
                    width: '200px',
                    height: '150px',
                    objectFit: 'cover',
                    display: 'block',
                  }}
                />
              </Box>
              <Box sx={{ mt: 1 }}>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowFileSelector(true)}
                >
                  Change Thumbnail
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setThumbnailId(null)}
                  sx={{ ml: 1 }}
                >
                  Remove
                </Button>
              </Box>
            </Box>
          ) : (
            <Button
              variant="outline"
              startIcon={<ImageIcon />}
              onClick={() => setShowFileSelector(true)}
            >
              Select Thumbnail
            </Button>
          )}
        </Box>

        <CategorySelector
          selectedCategories={selectedCategories}
          onChange={setSelectedCategories}
        />

        <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 3 }}>
          {onCancel && (
            <Button variant="outline" onClick={onCancel} disabled={loading}>
              Cancel
            </Button>
          )}
          <Button variant="fill" type="submit" disabled={loading}>
            {loading ? 'Saving...' : course ? 'Update Course' : 'Create Course'}
          </Button>
        </Box>
      </Paper>

      <FileSelector
        open={showFileSelector}
        bucketName="course-thumbnails"
        onClose={() => setShowFileSelector(false)}
        onSelect={handleThumbnailSelect}
        selectedFileId={thumbnailId}
        acceptedTypes={['image/*']}
      />
    </form>
  );
}

